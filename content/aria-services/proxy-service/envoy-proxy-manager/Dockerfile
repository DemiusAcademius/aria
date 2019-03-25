FROM debian:stretch-slim
WORKDIR /
ADD docker-content/bin/envoy-proxy-manager /
RUN chmod +x ./envoy-proxy-manager
CMD ["./envoy-proxy-manager"]