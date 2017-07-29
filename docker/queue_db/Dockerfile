FROM postgres:9.6

ENV POSTGRES_USER queue
COPY queue_db.sql docker-entrypoint-initdb.d
