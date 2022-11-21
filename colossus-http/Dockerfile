# Colossus Runner
#
# VERSION 0.0.1

FROM alpine:latest
LABEL Description="Runs Colossus binary"

RUN apk update && apk add curl

COPY . /opt/colossus/
WORKDIR /opt/colossus/

EXPOSE 9002

ENTRYPOINT ["./bin/colossus"]
