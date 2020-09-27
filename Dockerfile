FROM shineo/dev-game-data

FROM golang:1.15.2-alpine3.12

COPY build/wait-for /usr/local/bin

RUN chmod +x /usr/local/bin/wait-for

RUN mkdir /app

RUN apk add gcc

#ADD . /app

RUN mkdir /game-data

COPY --from=0 /game-data /game-data

WORKDIR /app

