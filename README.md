# Event Ticket Reservation System

Back-end

## Objectives:
- [Main] Learn and familiarize with Go programming language
- [Main] Concurrency in practice (with and w/o database features help)
- Pinpoint possible limitations in each implementation phase’s design
- Design for optimal performance (where to store data and in what form it should be stored)
## Notes:
- Event (Edit/Delete) not synced with Redis yet [Need to reflect event changes in Redis]
- Reservation (Cancel) not synced with Redis yet [Need to reclaim quotas after delete]

## Functional Requirements:
1. User System Functionalities
    - [x] Register / Sign-up
    - [x] Login
    - [x] User Roles
    - [x] Admin
    - [x] Event Organizer
    - [x] Customer
2. Event Management Functionalities
    - [x] Create/View/Edit/Delete Event
    * Required Properties
    - [x] Event Name
    - [x] Ticket Quota/Limit
- [x] Admins can view and delete all events in the system
- [x] Event Organizers can create/view/edit/delete their own events
- [x] Customers can view all events and are able to make reservation(s) on any event. (Customers can reserve more than 1 ticket per request.)
- [x] Tickets can only be sold if there’s quota left or have not reached their limits
- [x] Customers are able to cancel their reservations
- [x] Event Organizers can see total ticket reserved / remaining ticket quota

Design and implement a (backend) service with HTTP APIs serving the above mentioned functionalities using Go programming language

#### Phase 1 - Without using persistence (database). Store all data in the program's memory.
#### Phase 2 - Use persistence (relational database) (e.g. PostgreSQL) to store quota (naive implementation)
#### Phase 3 - Optimize and/or redesign for performance

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
| 4 Postgres + Redis| 2958 | 265.96 | 62724  |

## Docker useful commands
1. `docker build -t ticket_reservation_server .`
2. `docker run --rm -ti -p 9092:9092 ticket_reservation_server serve-api`
3. `docker run --rm -ti --entrypoint="/bin/bash" ticket_reservation_server`