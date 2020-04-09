# Inspired by: https://github.com/ory/hydra/blob/master/Dockerfile

# To compile this image manually run:
#
# $ GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o service && docker build -t shineo/service:local-build . && rm service
FROM alpine:3.11

RUN apk add -U --no-cache ca-certificates

FROM scratch

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY service /usr/bin/world

USER 1000

ENTRYPOINT ["world"]
CMD ["serve"]