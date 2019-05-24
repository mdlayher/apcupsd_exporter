FROM golang:alpine AS build-env

ADD . /src
RUN apk add --no-cache git
WORKDIR /src
RUN go build cmd/apcupsd_exporter/main.go

FROM alpine:latest

COPY --from=build-env /src/main /apcupsd_exporter

EXPOSE 9162

ENTRYPOINT ["/apcupsd_exporter"]
