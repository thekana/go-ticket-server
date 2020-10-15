# Ticket Reservation (Server)

## Prerequisites

- Go 1.13 or later

1. Start PostgreSQL

   _Example: Temporary PostgreSQL docker container_

   ```sh
   docker run --rm -p 5432:5432 -e POSTGRES_PASSWORD=postgres --name ticket_postgres postgres:12-alpine
   ```

   with mounting data dir from host

   ```sh
   docker run --rm -p 5432:5432 -e POSTGRES_PASSWORD=postgres -v $PWD/postgres_data:/var/lib/postgresql/data --name ticket_postgres postgres:12-alpine
   ```

2. Create `morchana_enterprise` database

   _Example: Temporary PostgreSQL docker container_

   ```sh
   docker exec -i ticket_postgres psql -U postgres -c "drop database if exists ticket_reservation" && \
   docker exec -i ticket_postgres psql -U postgres -c "create database ticket_reservation"
   ```

   _NOTE: Above command can also be used for resetting database_