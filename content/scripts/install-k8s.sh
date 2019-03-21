#!/bin/bash

BLUE='\033[1;34m'
RED='\033[1;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}INSTALL KUBERNETES SINGLE-NODE${NC}"

# timezone to Chisinau
timedatectl set-timezone Europe/Chisinau

# add kubernetes to packages registry
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

# install docker
echo -e "${BLUE}INSTALL DOCKER${NC}"
apt-get install docker.io -y
systemctl enable docker.service

# instal kubernetes utilities and services
echo -e "${BLUE}INSTALL K8S PACKAGES${NC}"
apt-get install -y kubelet kubeadm kubectl
apt-mark hold kubelet kubeadm kubectl

# preparing
systemctl daemon-reload
systemctl restart kubelet
sysctl net.bridge.bridge-nf-call-iptables=1

# run installation
echo -e "${BLUE}INSTALL KUBERNETES${NC}"
kubeadm init --pod-network-cidr=10.244.0.0/16

echo -e "${RED}OK${NC}"
