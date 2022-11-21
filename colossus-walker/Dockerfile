# Colossus Walker Runner
#
# VERSION 0.0.1

FROM alpine:latest
LABEL Description="Runs Colossus Walker binary"

RUN apk update && apk add curl

COPY . /opt/colossusw/
WORKDIR /opt/colossusw/

EXPOSE 9002

ENTRYPOINT ["./bin/colossusw"]
