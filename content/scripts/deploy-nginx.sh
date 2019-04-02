#!/bin/bash
VERSION=0.0.5

if $1 ; then
    kubectl delete service aria-nginx -n kube-system
    kubectl delete daemonset aria-nginx -n kube-system
fi

docker build -t aria-nginx:$VERSION aria-services/aria-nginx
docker tag aria-nginx:$VERSION 10.10.112.27:5000/aria-nginx:$VERSION
docker push 10.10.112.27:5000/aria-nginx:$VERSION

kubectl create -f aria-services/aria-nginx/daemonset.yaml
kubectl create -f aria-services/aria-nginx/service.yaml
