FROM envoyproxy/envoy-alpine:latest
RUN mkdir /var/log/envoy
COPY docker-content/envoy-config-managed.yaml /etc/envoy/envoy-config-managed.yaml
CMD /usr/local/bin/envoy -c /etc/envoy/envoy-config-managed.yaml