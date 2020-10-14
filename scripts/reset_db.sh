#!/bin/sh

docker exec -i ticket_reservation_postgres psql -U postgres -c "drop database if exists ticket_reservation" &&
docker exec -i ticket_reservation_postgres psql -U postgres -c "create database ticket_reservation"

go run main.go migrate-db
