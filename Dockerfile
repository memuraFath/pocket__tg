FROM golang:1.22.11-alpine3.21 AS builder

COPY . /github.com/memuraFath/pocket__tg/
WORKDIR /github.com/memuraFath/pocket__tg/

RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest
WORKDIR /root/

COPY --from=0 /github.com/memuraFath/pocket__tg/bin/bot .
COPY --from=0 /github.com/memuraFath/pocket__tg/configs configs

EXPOSE 80

CMD ["./bot"]
