FROM golang:1.25.3 AS builder

WORKDIR /app


COPY . .

RUN go build -o manage .

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /app/manage /app/manage
EXPOSE 8080

CMD ["/app/manage"]
