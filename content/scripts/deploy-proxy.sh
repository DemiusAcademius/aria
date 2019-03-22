#!/bin/bash
VERSION=0.0.11
kubectl delete daemonset proxy-service -n kube-system

docker build -t envoy-proxy.aria:$VERSION aria-services/proxy-service/envoy-proxy
docker tag envoy-proxy.aria:$VERSION 10.10.112.27:5000/envoy-proxy.aria:$VERSION
docker push 10.10.112.27:5000/envoy-proxy.aria:$VERSION

docker build -t envoy-proxy-manager.aria:$VERSION aria-services/proxy-service/envoy-proxy-manager
docker tag envoy-proxy-manager.aria:$VERSION 10.10.112.27:5000/envoy-proxy-manager.aria:$VERSION
docker push 10.10.112.27:5000/envoy-proxy-manager.aria:$VERSION

kubectl create -f aria-services/proxy-service/daemonset.yaml
