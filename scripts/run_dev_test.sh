#!/bin/bash

trap killgroup SIGINT

killgroup(){
  echo killing...
  kill 0
}

if [ ! "$(docker ps -q -f name=ticket_reservation_postgres)" ]; then
    if [ "$(docker ps -aq -f status=exited -f name=ticket_reservation_postgres)" ]; then
        # cleanup
        docker rm ticket_reservation_postgres
    fi
    # run your container
    docker run -d --rm -p 5432:5432 -e POSTGRES_PASSWORD=postgres --name ticket_reservation_postgres postgres:12-alpine
    sleep 2
fi

if [ ! "$(docker ps -q -f name=ticket_reservation_redis)" ]; then
    if [ "$(docker ps -aq -f status=exited -f name=ticket_reservation_redis)" ]; then
        # cleanup
        docker rm ticket_reservation_redis
    fi
    # run your container
    docker run -d --rm -p 6379:6379 --name ticket_reservation_redis redis:6-alpine
fi

docker exec -i ticket_reservation_postgres psql -U postgres -c "drop database if exists ticket_reservation" &&
docker exec -i ticket_reservation_postgres psql -U postgres -c "create database ticket_reservation"
docker exec -i ticket_reservation_redis redis-cli FLUSHALL

go run main.go migrate-db
go run main.go seed-db
go run -tags debug main.go serve-api &

wait