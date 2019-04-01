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

# connect k8s to registry
kubectl create secret docker-registry regcred --docker-server=10.10.112.27:5000 --docker-username=kube-registry-user --docker-password=yf-ujhirt-cbltk-rjhjkm --docker-email=demius.md@gmail.com -n kube-system
kubectl create secret docker-registry regcred --docker-server=10.10.112.27:5000 --docker-username=kube-registry-user --docker-password=yf-ujhirt-cbltk-rjhjkm --docker-email=demius.md@gmail.com -n admin
kubectl create secret docker-registry regcred --docker-server=10.10.112.27:5000 --docker-username=kube-registry-user --docker-password=yf-ujhirt-cbltk-rjhjkm --docker-email=demius.md@gmail.com -n billing
kubectl create secret docker-registry regcred --docker-server=10.10.112.27:5000 --docker-username=kube-registry-user --docker-password=yf-ujhirt-cbltk-rjhjkm --docker-email=demius.md@gmail.com -n finance
kubectl create secret docker-registry regcred --docker-server=10.10.112.27:5000 --docker-username=kube-registry-user --docker-password=yf-ujhirt-cbltk-rjhjkm --docker-email=demius.md@gmail.com -n utilities

# deploy proxy-service
echo ""
echo -e "${BLUE}Deploy PROXY${NC}"
PROXY_VERSION=0.1.2
docker build -t aria-proxy-service:$PROXY_VERSION aria-services/aria-proxy/aria-proxy-service
docker tag aria-proxy-service:$PROXY_VERSION 10.10.112.27:5000/aria-proxy-service:$PROXY_VERSION
docker push 10.10.112.27:5000/aria-proxy-service:$PROXY_VERSION

docker build -t aria-proxy-manager:$PROXY_VERSION aria-services/aria-proxy/aria-proxy-manager
docker tag aria-proxy-manager:$PROXY_VERSION 10.10.112.27:5000/aria-proxy-manager:$PROXY_VERSION
docker push 10.10.112.27:5000/aria-proxy-manager:$PROXY_VERSION

kubectl create -f aria-services/aria-proxy/daemonset.yaml


# deploy nginx
echo ""
echo -e "${BLUE}Deploy static web-server${NC}"
NGINX_VERSION=0.0.5
docker build -t aria-nginx:$NGINX_VERSION aria-services/aria-nginx
docker tag aria-nginx:$NGINX_VERSION 10.10.112.27:5000/aria-nginx:$NGINX_VERSION
docker push 10.10.112.27:5000/aria-nginx:$NGINX_VERSION

kubectl create -f aria-services/aria-nginx/daemonset.yaml
kubectl create -f aria-services/aria-nginx/service.yaml

# find bearer token for login to dashboard
echo ""
echo -e "${RED}COPY DASHBOARD AUTH TOKEN:${NC}"
kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep dashboard-user | awk '{print $1}')

chmod +x scripts/sync-applications

echo -e "${RED}OK${NC}"
echo -e "${RED}run scripts/sync-applications${NC}"
