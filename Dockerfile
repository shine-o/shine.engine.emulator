FROM shineo/dev-game-data

FROM golang:1.15.2-alpine3.12

RUN mkdir /app

#ADD . /app

RUN mkdir /game-data

COPY --from=0 /game-data /game-data

WORKDIR /app

