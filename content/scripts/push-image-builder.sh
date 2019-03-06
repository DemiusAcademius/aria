#!/bin/bash

microk8s.docker build -t localhost:32000/image-builder.aria:0.0.1 applications/System/image-builder
microk8s.docker push localhost:32000/image-builder.aria:0.0.1