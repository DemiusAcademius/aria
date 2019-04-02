#!/bin/bash
VERSION=0.1.1

if $1 ; then
    kubectl delete daemonset aria-publisher -n kube-system
fi

docker build -t aria-publisher:$VERSION aria-services/aria-publisher
docker tag aria-publisher:$VERSION 10.10.112.27:5000/aria-publisher:$VERSION
docker push 10.10.112.27:5000/aria-publisher:$VERSION

kubectl create -f aria-services/aria-publisher/daemonset.yaml
