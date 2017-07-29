FROM postgres:9.6

RUN mkdir -p /healthcheck
COPY waitforit /healthcheck/
HEALTHCHECK --interval=15s --timeout=5s --retries=5 CMD /healthcheck/waitforit --host=localhost --port=5432 --timeout=5

ENV POSTGRES_USER queue
COPY queue_db.sql docker-entrypoint-initdb.d
