FROM openjdk:11.0.3-jre

VOLUME /tmp
ADD {{.Executable}}-{{.Version}}.jar app.jar

ENTRYPOINT ["java","-Dspring.profiles.active=production","-Dfile.encoding=UTF-8","-Djava.security.egd=file:/dev/./urandom","-jar","/app.jar"]
