FROM golang:1.14.10-alpine3.12
WORKDIR /ticket-reservation-server
COPY . .
RUN go build -o main .
EXPOSE 9092
ENTRYPOINT ["/ticket-reservation-server/main"]