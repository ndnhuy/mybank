FROM eclipse-temurin:21-jdk-alpine as build

WORKDIR /app

COPY .mvn/ .mvn
COPY mvnw pom.xml ./
COPY .m2 /root/.m2
# RUN ./mvnw dependency:go-offline
COPY src ./src
RUN  ./mvnw install -DskipTests

FROM eclipse-temurin:21-jre-alpine
COPY --from=build /app/target/*.jar app.jar
ENTRYPOINT [ "java", "-jar", "-Dspring.profiles.active=default,docker", "-Dcom.sun.management.jmxremote=true", "-Djava.rmi.server.hostname=127.0.0.1", "-Dcom.sun.management.jmxremote.host=0.0.0.0", "-Dcom.sun.management.jmxremote.port=9991", "-Dcom.sun.management.jmxremote.rmi.port=9991", "-Dcom.sun.management.jmxremote.ssl=false", "-Dcom.sun.management.jmxremote.registry.ssl=false", "-Dcom.sun.management.jmxremote.authenticate=false", "-Djava.net.preferIPv4Stack=true", "/app.jar" ]