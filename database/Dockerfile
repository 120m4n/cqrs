FROM postgres:12.1

# Copy the database initialization script:
COPY up.sql /docker-entrypoint-initdb.d/1.sql


CMD ["postgres"]