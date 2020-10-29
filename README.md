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

3. Create tables and run server

   ```sh
   go run main.go migrate-db
   go run -tags debug main.go serve-api
   ```
4. Connect to psql
    ```shell script
    docker exec -ti ticket_reservation_postgres psql -U postgres
    ```
## Performance Test

Uses autocannon with config:
1. `connections: 80`
2. `duration: 10`
3. `4 POST requests. Each request reserves 1 ticket from 1 event`
4. `excludeErrorStats: true`

Start by sending a GET request to populate DB (ONCE) with
- EventID 1 Quota 10000
- EventID 2 Quota 10000
- EventID 3 Quota 10000
- EventID 4 Quota 10000
```shell script
curl localhost:9092/api/v1/pop
```
Then run `load_test/cannon.js` with updated customer tokens <br>
P.S. Need tokens because my implementation uses jwtclaims to get user data 🤣

## Results
| Phase    | Requests Sent | Avg Latency     | Avg Throughput     |
| :-------------| :----------: | :----------: | -----------: |
|  1 Memory only| 9307 | 85.97   | 300800   |
| 2 Postgres only (with retry) | 713 | 513.71 |10302.21  |
| 3 Postgres with batch jobs in memory| 4992 | 161.47 | 119737.6  |