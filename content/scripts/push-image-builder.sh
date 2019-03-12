#!/bin/bash

docker build -t image-builder.aria:0.0.6 aria-services/image-builder
docker tag image-builder.aria:0.0.6 10.10.112.27:5000/image-builder.aria:0.0.6
docker push 10.10.112.27:5000/image-builder.aria:0.0.6