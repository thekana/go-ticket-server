#!/bin/sh

docker run --rm -p 5432:5432 -e POSTGRES_PASSWORD=postgres --name ticket_reservation_postgres postgres:12-alpine
