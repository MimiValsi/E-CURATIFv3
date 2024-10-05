FROM postgres:16-alpine

ENV POSTGRES_USER=postgres

COPY create_db.sql /docker-entrypoint-initdb.d/
