FROM golang:1.22-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/hyle-team/bridgeless-signer

ENV GO111MODULE="on"
ENV CGO_ENABLED=1
ENV GOOS="linux"
ENV GOPRIVATE=github.com/*
ENV GONOSUMDB=github.com/*
ENV GONOPROXY=github.com/*

COPY ./go.mod ./go.sum ./
# Read the CI_ACCESS_TOKEN from the .env file
ARG CI_ACCESS_TOKEN
RUN git config --global url."https://olegfomenkodev:${CI_ACCESS_TOKEN}@github.com/".insteadOf "https://github.com/"
RUN go mod download

COPY . .

RUN go mod vendor
RUN go build  -o /usr/local/bin/bridgeless-signer /go/src/github.com/hyle-team/bridgeless-signer


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/bridgeless-signer /usr/local/bin/bridgeless-signer

RUN apk add --no-cache ca-certificates

ENTRYPOINT ["bridgeless-signer"]
