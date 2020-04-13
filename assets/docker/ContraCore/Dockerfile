# Build

FROM golang:1.14.2 AS build-dev

RUN mkdir -p /opt/

COPY build.sh /opt/
COPY Corefile /opt/

RUN /bin/bash /opt/build.sh

# Run

FROM ubuntu:eoan

COPY --from=build-dev /opt/coredns/coredns /opt/coredns
COPY Corefile /opt/
COPY wait-for-it.sh /opt/wait-for-it.sh

EXPOSE 53/tcp
EXPOSE 53/udp

CMD /opt/coredns -conf /opt/Corefile
