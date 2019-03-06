#!/bin/bash

docker build -t image-builder.aria:0.0.1 applications/System/image-builder
docker tag image-builder.aria:0.0.1 127.0.0.1:32000/image-builder.aria:0.0.1
docker push 127.0.0.1:32000/image-builder.aria:0.0.1