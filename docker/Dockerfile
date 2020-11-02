FROM golang:1.14-alpine3.11 as builder
WORKDIR /ticket-reservation-server/
COPY . .
RUN go build

FROM alpine:3.11
RUN apk --no-cache add ca-certificates

# Install bash shell for convenience.
RUN apk add --no-cache bash
#COPY --from=builder /ticket-reservation-server/_dev_server_keys ./_dev_server_keys
#COPY --from=builder /ticket-reservation-server/config.yaml .
COPY --from=builder /ticket-reservation-server/ticket-reservation .
COPY ./docker/wait-for-it.sh .

ENTRYPOINT ["/ticket-reservation"]