FROM debian:stretch-slim
WORKDIR /
ADD docker-content/bin/aria-proxy-manager /
RUN chmod +x ./aria-proxy-manager
CMD ["./aria-proxy-manager"]