#!/bin/bash

BLUE='\033[1;36m'
RED='\033[1;31m'
NC='\033[0m' # No Color

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
systemctl daemon-reload
systemctl restart kubelet
sysctl net.bridge.bridge-nf-call-iptables=1

# run installation
echo ""
echo -e "${BLUE}INSTALL KUBERNETES${NC}"
kubeadm init --pod-network-cidr=10.244.0.0/16

echo -e "${RED}  COPY AND SAVE THE ABOVE!${NC}"
echo ""

# post install
mkdir -p $HOME/.kube
cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
chown -R acc-server-admin:acc-server-admin $HOME/.kube/config

# setup resolv.conf in pods
cp assets/resolv.conf /etc/kubernetes/resolv.conf
sed -i 's:/run/systemd/resolve/resolv.conf:/etc/kubernetes/resolv.conf:g' /var/lib/kubelet/kubeadm-flags.env
systemctl restart kubelet

# tuning core-dns config
kubectl delete cm coredns -n kube-system
kubectl create -f manifests/dns-config.yaml
kubectl delete pod -l k8s-app=kube-dns -n kube-system

# install pod network addon (flannel)
# see: https://kubernetes.io/docs/setup/independent/create-cluster-kubeadm/
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/a70459be0084506e4ec919aa1c114638878db11b/Documentation/kube-flannel.yml

echo -e "${RED}OK${NC}"
echo -e "${RED}SYSTEM WILL REBOOT NOW${NC}"
echo -e "${RED}run scripts/install-after-reboot${NC}"
reboot

