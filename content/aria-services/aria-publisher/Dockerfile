FROM debian:stretch-slim
WORKDIR /
ADD docker-content/bin/aria-publisher /
RUN chmod +x ./aria-publisher
CMD ["./aria-publisher"]