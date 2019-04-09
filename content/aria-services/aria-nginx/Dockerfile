FROM docker.io/nginx:alpine

COPY docker-content/default.conf /etc/nginx/conf.d

COPY docker-content/index.html /usr/share/nginx/html
ADD docker-content/media.tar.gz /usr/share/nginx/html/media

RUN chmod -R +r /usr/share/nginx/html/media
