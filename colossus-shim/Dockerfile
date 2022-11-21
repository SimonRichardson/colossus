# Colossus Shim Runner
#
# VERSION 0.0.1

FROM alpine:latest
LABEL Description="Runs Colossus Shim binary"

RUN apk update && apk add curl

COPY . /opt/colossuss/
WORKDIR /opt/colossuss/

EXPOSE 9202

ENTRYPOINT ["./bin/colossuss"]
