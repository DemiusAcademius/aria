#!/bin/bash
VERSION=0.1.0
kubectl delete daemonset aria-publisher -n kube-system

docker build -t aria-publisher:$VERSION aria-services/aria-publisher
docker tag aria-publisher:$VERSION 10.10.112.27:5000/aria-publisher:$VERSION
docker push 10.10.112.27:5000/aria-publisher:$VERSION

kubectl create -f aria-services/aria-publisher/daemonset.yaml
