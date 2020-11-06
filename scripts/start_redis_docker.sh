#!/bin/sh

docker run --rm -p 6379:6379 --name ticket_reservation_redis redis:6-alpine
