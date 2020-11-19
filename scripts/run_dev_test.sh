#!/bin/bash

trap killgroup SIGINT

killgroup(){
  echo killing...
  kill 0
}

# docker run --rm -p 5432:5432 -e POSTGRES_PASSWORD=postgres --name ticket_reservation_postgres postgres:10-alpine &
# sleep 3

docker exec -i ticket_reservation_postgres psql -U postgres -c "drop database if exists ticket_reservation" &&
docker exec -i ticket_reservation_postgres psql -U postgres -c "create database ticket_reservation"

go run main.go migrate-db
go run main.go seed-db
go run -tags debug main.go serve-api &

wait