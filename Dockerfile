FROM golang:1.16-alpine

RUN apk add gcc

RUN mkdir /app

ADD . /app

RUN chmod +x /app/build/wait-for

WORKDIR /app

RUN go mod download