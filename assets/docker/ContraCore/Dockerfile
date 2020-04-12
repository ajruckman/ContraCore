FROM golang:1.14.2

RUN mkdir -p /opt/

COPY build.sh /opt/

RUN /bin/bash /opt/build.sh

CMD tail -f /dev/null
