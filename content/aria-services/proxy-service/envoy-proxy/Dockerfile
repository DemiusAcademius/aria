FROM envoyproxy/envoy-alpine:latest
RUN mkdir /var/log/envoy
COPY docker-content/envoy-config.yaml /etc/envoy/envoy-config.yaml
CMD /usr/local/bin/envoy -c /etc/envoy/envoy-config.yaml