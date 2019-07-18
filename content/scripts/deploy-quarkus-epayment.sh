#!/bin/bash
VERSION=0.1

if $1 ; then
    kubectl delete service quarkus-epayment-api -n billing
    kubectl delete deployment quarkus-epayment.api -n billing
fi

docker build -t quarkus-epayment.api.billing:$VERSION aria-services/quarkus-epayment
docker tag quarkus-epayment.api.billing:$VERSION 10.10.112.27:5000/quarkus-epayment.api.billing:$VERSION
docker push 10.10.112.27:5000/quarkus-epayment.api.billing:$VERSION

kubectl create -f aria-services/quarkus-epayment/deployment.yaml
kubectl create -f aria-services/quarkus-epayment/service.yaml
