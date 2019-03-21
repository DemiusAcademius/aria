#!/bin/bash

BLUE='\033[1;34m'
RED='\033[1;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}POST INSTALL${NC}"

# Create dashboard, dashboard user && namespaces
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml
kubectl create -f manifests/dashboard-user.yaml
kubectl create -f manifests/namespaces.yaml

# copy certificates
cp certs/*.crt /usr/local/share/ca-certificates
update-ca-certificates
systemctl restart docker.service

# install docker registry
echo -e "${BLUE}Install registry${NC}"
mkdir registry-data
docker run --entrypoint htpasswd registry:2 -Bbn kube-registry-user yf-ujhirt-cbltk-rjhjkm > ./auth/htpasswd
chown root:root ./auth/htpasswd
kubectl create -f manifests/docker-registry-daemonset.yaml

# Connect k8s to registry
kubectl create secret docker-registry regcred --docker-server=10.10.112.27:5000 --docker-username=kube-registry-user --docker-password=yf-ujhirt-cbltk-rjhjkm --docker-email=demius.md@gmail.com

echo -e "${RED}OK${NC}"
