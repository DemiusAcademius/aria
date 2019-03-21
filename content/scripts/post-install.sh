#!/bin/bash

BLUE='\033[1;36m'
RED='\033[1;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}POST INSTALL${NC}"

# master isolation for single node installation
kubectl taint nodes --all node-role.kubernetes.io/master-

# create dashboard, dashboard user && namespaces
# SEE: https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml
echo -e "${BLUE}Install Dashboard${NC}"
kubectl create secret generic kubernetes-dashboard-certs --from-file=$HOME/certs/dashboard -n kube-system
kubectl create -f manifests/dashboard-service.yaml
kubectl create -f manifests/dashboard-user.yaml
kubectl create -f manifests/namespaces.yaml

# install docker registry
echo -e "${BLUE}Install registry${NC}"
mkdir registry-data
docker run --entrypoint htpasswd registry:2 -Bbn kube-registry-user yf-ujhirt-cbltk-rjhjkm > ./auth/htpasswd
chown root:root ./auth/htpasswd
kubectl create -f manifests/docker-registry-daemonset.yaml

# connect k8s to registry
kubectl create secret docker-registry regcred --docker-server=10.10.112.27:5000 --docker-username=kube-registry-user --docker-password=yf-ujhirt-cbltk-rjhjkm --docker-email=demius.md@gmail.com

# find bearer token for login to dashboard
echo ""
echo -e "${RED}COPY DASHBOARD AUTH TOKEN:${NC}"
kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep dashboard-user | awk '{print $1}')

echo -e "${RED}OK${NC}"