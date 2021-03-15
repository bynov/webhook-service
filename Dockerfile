FROM golang:1.16.2-buster

WORKDIR /app

COPY . .

RUN go build

FROM debian:buster

WORKDIR /app

COPY --from=0 /app/webhook-service .

CMD ["./webhook-service"]
