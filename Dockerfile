FROM golang:1.18.3-alpine3.16

WORKDIR /usr/src/app

COPY . .
RUN go mod download && go mod verify
RUN go build -v -o main

ARG TELEGRAM_APITOKEN=prod

ENV TELEGRAM_APITOKEN $TELEGRAM_APITOKEN

CMD ["/usr/src/app/main"]