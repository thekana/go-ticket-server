# Ticket Reservation (Server)

## Prerequisites

- Go 1.13 or later

1. Start PostgreSQL

   ```sh
   docker run --rm -p 5432:5432 -e POSTGRES_PASSWORD=postgres --name ticket_reservation_postgres postgres:12-alpine
   ```

2. Create database

   ```sh
   docker exec -i ticket_reservation_postgres psql -U postgres -c "drop database if exists ticket_reservation" &&
   docker exec -i ticket_reservation_postgres psql -U postgres -c "create database ticket_reservation"
   ```

3. Create tables and server

   ```sh
   go run main.go migrate-db
   go run -tags debug main.go serve-api
   ```