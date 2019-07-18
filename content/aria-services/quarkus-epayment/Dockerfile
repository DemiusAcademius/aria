FROM quarkus-docker-base

VOLUME /tmp
ADD docker-content/quarkus-epayment-api-runner.jar app.jar
ADD docker-content/certs.tar.gz /certs

ENTRYPOINT ["java","-Dfile.encoding=UTF-8","-Djava.security.egd=file:/dev/./urandom","-jar","/app.jar"]