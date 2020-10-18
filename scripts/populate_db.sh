#!/bin/sh
cat ./scripts/query.sql | docker exec -i ticket_reservation_postgres psql -U postgres -d ticket_reservation