FROM golang:1.25.1 AS builder

WORKDIR /app


COPY . .

RUN go build -o electronic .

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/electronic /app/electronic
COPY --from=builder /app/pkg/dbwork/migrations /app/pkg/dbwork/migrations
EXPOSE 8080

CMD ["/app/electronic"]

