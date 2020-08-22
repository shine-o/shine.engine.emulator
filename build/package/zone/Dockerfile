# Inspired by: https://github.com/ory/hydra/blob/master/Dockerfile

# To compile this image manually run:
#
# $ GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o world && docker build -t shineo/service:local-build . && rm service
FROM alpine:3.11

RUN apk add -U --no-cache ca-certificates

FROM scratch

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY zone /usr/bin/zone

USER 1000

ENTRYPOINT ["zone"]
CMD ["serve"]