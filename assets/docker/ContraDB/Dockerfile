FROM postgres:12

#RUN mkdir -p /opt/contradb

COPY *.sql /docker-entrypoint-initdb.d/
#COPY contradb.sh /docker-entrypoint-initdb.d/
#COPY contradb/* /opt/contradb/

ENV POSTGRES_PASSWORD TNwNmkyaiM29w4gR

EXPOSE 5432/tcp
