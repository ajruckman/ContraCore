FROM yandex/clickhouse-server

COPY *.sql /docker-entrypoint-initdb.d/
COPY contralog_mgr.xml /etc/clickhouse-server/users.d/

EXPOSE 8123/tcp
EXPOSE 9000/tcp
