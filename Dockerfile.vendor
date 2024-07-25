FROM golang:1.22-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/hyle-team/bridgeless-signer
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/bridgeless-signer /go/src/github.com/hyle-team/bridgeless-signer


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/bridgeless-signer /usr/local/bin/bridgeless-signer
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["bridgeless-signer"]
