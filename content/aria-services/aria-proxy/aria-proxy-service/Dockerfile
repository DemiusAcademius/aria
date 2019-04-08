FROM envoyproxy/envoy-alpine:v1.10.0
RUN mkdir /var/log/envoy
COPY docker-content/envoy-config-managed.yaml /etc/envoy/envoy-config-managed.yaml
CMD /usr/local/bin/envoy -c /etc/envoy/envoy-config-managed.yaml