#!/bin/bash
VERSION=0.0.3
# kubectl delete service nginx -n kube-system
kubectl delete daemonset nginx -n kube-system

docker build -t nginx.aria:$VERSION aria-services/nginx
docker tag nginx.aria:$VERSION 10.10.112.27:5000/nginx.aria:$VERSION
docker push 10.10.112.27:5000/nginx.aria:$VERSION

kubectl create -f aria-services/nginx/daemonset.yaml
# kubectl create -f aria-services/nginx/service.yaml