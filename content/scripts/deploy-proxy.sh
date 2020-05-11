#!/bin/bash
VERSION=0.1.7

if $1 ; then
    kubectl delete daemonset aria-proxy -n kube-system
fi

docker build -t aria-proxy-service:$VERSION aria-services/aria-proxy/aria-proxy-service
docker tag aria-proxy-service:$VERSION 10.10.112.27:5000/aria-proxy-service:$VERSION
docker push 10.10.112.27:5000/aria-proxy-service:$VERSION

docker build -t aria-proxy-manager:$VERSION aria-services/aria-proxy/aria-proxy-manager
docker tag aria-proxy-manager:$VERSION 10.10.112.27:5000/aria-proxy-manager:$VERSION
docker push 10.10.112.27:5000/aria-proxy-manager:$VERSION

kubectl create -f aria-services/aria-proxy/daemonset.yaml
