#!/bin/bash
VERSION=0.0.9
kubectl delete daemonset image-builder -n aria

docker build -t image-builder.aria:$VERSION aria-services/image-builder
docker tag image-builder.aria:$VERSION 10.10.112.27:5000/image-builder.aria:$VERSION
docker push 10.10.112.27:5000/image-builder.aria:$VERSION

kubectl create -f aria-services/image-builder/daemonset.yaml
## kubectl create -f aria-services/image-builder/service.yaml