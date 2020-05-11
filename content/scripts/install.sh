#!/bin/bash

BLUE='\033[1;36m'
RED='\033[1;31m'
NC='\033[0m' # No Color
ADMIN_USER_HOME=/home/acc-server-admin

echo -e "${BLUE}INSTALL KUBERNETES SINGLE-NODE${NC}"
echo ""

# timezone to Chisinau
timedatectl set-timezone Europe/Chisinau

# config /etc/hosts
echo '127.0.0.1   kubernetes-srv' >> /etc/hosts

# add kubernetes to packages registry
echo -e "${BLUE}add kubernetes packages to registry${NC}"
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/ kubernetes-xenial main
EOF

# update packages registry & update
echo -e "${BLUE}UPDATE SYSTEM${NC}"
apt update & apt upgrade -y

# disable swap
sed -ri '/\sswap\s/s/^#?/#/' /etc/fstab
swapoff -a

# copy certificates
cp certs/acc.io/*.crt /usr/local/share/ca-certificates
cp certs/acc.md/*.crt /usr/local/share/ca-certificates
update-ca-certificates

# install docker
echo ""
echo -e "${BLUE}INSTALL DOCKER${NC}"
apt-get install docker.io -y
systemctl enable docker.service

# instal kubernetes utilities and services
echo ""
echo -e "${BLUE}INSTALL K8S PACKAGES${NC}"
apt-get install -y kubelet kubeadm kubectl
apt-mark hold kubelet kubeadm kubectl

# preparing

cat <<EOF > /etc/sysctl.d/99-kubernetes-cri.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.ip_forward=1
net.ipv4.conf.all.rp_filter=1
net.netfilter.nf_conntrack_max = 1000000
EOF
sysctl --system

systemctl daemon-reload
systemctl restart kubelet

# Prevent NetworkManager from interfering with the interfaces:
cat <<EOF > /etc/NetworkManager/conf.d/calico.conf
[keyfile]
unmanaged-devices=interface-name:cali*;interface-name:tunl*
EOF
systemctl restart NetworkManager  > /dev/null 2>&1

# run installation
echo ""
echo -e "${BLUE}INSTALL KUBERNETES${NC}"
kubeadm init --pod-network-cidr=192.168.0.0/16

echo -e "${RED}  COPY AND SAVE THE ABOVE!${NC}"
echo ""

# post install
# for admin user
mkdir -p $ADMIN_USER_HOME/.kube
cp -i /etc/kubernetes/admin.conf $ADMIN_USER_HOME/.kube/config
chown -R acc-server-admin:acc-server-admin $ADMIN_USER_HOME/.kube/config

# for root
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config

# setup resolv.conf in pods
# cp assets/resolv.conf /etc/kubernetes/resolv.conf
# sed -i 's:/run/systemd/resolve/resolv.conf:/etc/kubernetes/resolv.conf:g' /var/lib/kubelet/kubeadm-flags.env

systemctl restart kubelet
sleep 2s

# tuning core-dns config
# kubectl delete cm coredns -n kube-system
# kubectl create -f manifests/dns-config.yaml
# kubectl delete pod -l k8s-app=kube-dns -n kube-system

# install pod network addon (Calico)
# see: https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/
kubectl apply -f https://docs.projectcalico.org/v3.11/manifests/calico.yaml

echo -e "${RED}OK${NC}"
echo -e "${RED}SYSTEM WILL REBOOT NOW${NC}"
echo -e "${RED}run sudo scripts/install-after-reboot.sh${NC}"
reboot

