FROM postgres:14

COPY init-user-db.sh /docker-entrypoint-initdb.d/01.sh
COPY schema.sql /tmp/schema.sql