# as sudo

# install docker etc
dnf update -y
dnf install docker nano wget -y

# Ensure iptables tooling does not use the nftables backend
update-alternatives --set iptables /usr/sbin/iptables-legacy

# Make sure that the br_netfilter module is loaded
cat > /etc/modules-load.d/containerd.conf <<EOF
overlay
br_netfilter
EOF

modprobe overlay
modprobe br_netfilter

# Setup required sysctl params, these persist across reboots.
cat > /etc/sysctl.d/99-kubernetes-cri.conf <<EOF
net.bridge.bridge-nf-call-iptables  = 1
net.ipv4.ip_forward                 = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.conf.all.rp_filter         = 1
EOF

sysctl --system

# start docker
systemctl enable --now docker

# Installing kubeadm, kubelet and kubectl
cat <<EOF > /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-x86_64
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF

# Set SELinux in permissive mode (effectively disabling it)
setenforce 0
sed -i 's/^SELINUX=enforcing$/SELINUX=permissive/' /etc/selinux/config

# install kubelet etc
dnf install -y kubelet kubeadm kubectl --disableexcludes=kubernetes

# start kubelet
systemctl enable --now kubelet

# config firewall
firewall-cmd --permanent --add-port=6443/tcp
firewall-cmd --permanent --add-port=6443/udp
firewall-cmd --permanent --add-port=10250/tcp
firewall-cmd --permanent --add-port=10250/udp
firewall-cmd --permanent --add-port=30931/tcp
firewall-cmd --reload

# create a single-host kubernetes cluster with Calico
sudo kubeadm init --pod-network-cidr=192.168.0.0/16

# exit from su

# config admin user
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# install dashboard
wget https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta6/aio/deploy/recommended.yaml
mv recommended.yaml dashboard.yaml
nano dashboard.yaml

  ports:
    - port: 443
      targetPort: 8443
      nodePort: 30931   <---
  selector:
    k8s-app: kubernetes-dashboard
  type: NodePort        <---

kubectl apply -f dashboard.yaml

# setup coredns resolv
nano /etc/resolv.conf

nameserver 10.96.0.10
# nameserver 10.10.0.250
# nameserver 10.10.0.252
search namespace.svc.cluster.local svc.cluster.local cluster.local apa-canal.md
options ndots:5

nano /var/lib/kubelet/kubeadm-flags.env
# add  --resolv-conf=/etc/kubernetes/resolv.conf


