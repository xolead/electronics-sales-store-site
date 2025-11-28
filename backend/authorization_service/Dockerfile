FROM golang:1.25.3 AS builder

WORKDIR /app


COPY . .

RUN go build -o authoriz .

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/auth /app/authoriz
COPY --from=builder /app/pkg/dbwork/migrations /app/pkg/dbwork/migrations
EXPOSE 8083

CMD ["/app/authoriz"]
