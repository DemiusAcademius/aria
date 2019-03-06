FROM debian:stretch-slim

RUN rm -rf /var/lib/apt/lists/*

WORKDIR /
ADD docker-content/image-builder /
RUN chmod +x ./image-builder
CMD ["./image-builder"]