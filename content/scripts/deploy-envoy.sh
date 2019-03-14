#!/bin/bash
VERSION=0.0.5
kubectl delete daemonset envoy-proxy -n aria

docker build -t envoy-proxy.aria:$VERSION aria-services/envoy-proxy
docker tag envoy-proxy.aria:$VERSION 10.10.112.27:5000/envoy-proxy.aria:$VERSION
docker push 10.10.112.27:5000/envoy-proxy.aria:$VERSION

kubectl create -f aria-services/envoy-proxy/daemonset.yaml
