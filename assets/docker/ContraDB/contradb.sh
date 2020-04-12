#!/bin/bash

url="postgresql://contracore_mgr:uTiXe3oYJDv9Z4Ef@localhost/contradb"

# https://stackoverflow.com/a/39028690/9911189
retries=5
until psql $url -c 'select 1' >/dev/null 2>&1 || [ $retries -eq 0 ]; do
  echo "Waiting for postgres server, $((retries--)) remaining attempts..."
  sleep 1
done

psql $url -f /opt/contradb/10_ddl.sql
psql $url -f /opt/contradb/20_stats.sql
