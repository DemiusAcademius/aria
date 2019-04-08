#!/bin/bash

BLUE='\033[1;36m'
RED='\033[1;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}POST INSTALL${NC}"
echo ""

# master isolation for single node installation
kubectl taint nodes --all node-role.kubernetes.io/master-

# create dashboard, dashboard user && namespaces
# SEE: https://raw.githubusercontent.com/kubernetes/dashboard/v1.10.1/src/deploy/recommended/kubernetes-dashboard.yaml
echo ""
echo -e "${BLUE}Install Dashboard${NC}"
kubectl create secret generic kubernetes-dashboard-certs --from-file=$HOME/certs/dashboard -n kube-system
kubectl create -f manifests/dashboard-service.yaml
kubectl create -f manifests/dashboard-user.yaml
kubectl create -f manifests/namespaces.yaml
kubectl create -f manifests/aria-services-account.yaml

# install docker registry
echo ""
echo -e "${BLUE}Install registry${NC}"
mkdir registry-data
docker run --entrypoint htpasswd registry:2 -Bbn kube-registry-user yf-ujhirt-cbltk-rjhjkm > ./auth/htpasswd
chown root:root ./auth/htpasswd
kubectl create -f manifests/docker-registry.yaml

# connect k8s to registry and create configmaps
array=( kube-system admin billing finance utilities )
for i in "${array[@]}"
do
    kubectl create secret generic acc-certs --from-file=$HOME/certs/acc.io -n $i
    kubectl create secret generic auth-token-secret --from-file=$HOME/assets/asp.net/auth-token-secret -n $i
    kubectl create secret docker-registry regcred --docker-server=10.10.112.27:5000 --docker-username=kube-registry-user --docker-password=yf-ujhirt-cbltk-rjhjkm --docker-email=demius.md@gmail.com -n $i
    kubectl create configmap nginx-conf --from-file=assets/nginx/default.conf -n $i
done

# find bearer token for login to dashboard
echo ""
echo -e "${RED}COPY DASHBOARD AUTH TOKEN:${NC}"
kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep dashboard-user | awk '{print $1}')

chmod +x scripts/sync-applications

echo -e "${RED}OK${NC}"
echo -e "${RED}run sudo scripts/install-aria.sh${NC}"
