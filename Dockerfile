FROM golang:1.12 as build

ARG GO111MODULE=on
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

WORKDIR /go/src/app
COPY . ./

RUN go build \
    -a \
    -tags netgo \
    -ldflags '-w -extldflags "-static"' \
    -o apcupsd_exporter \
    cmd/apcupsd_exporter/main.go

FROM scratch
COPY --from=build /go/src/app/apcupsd_exporter /

EXPOSE 9162
ENTRYPOINT [ "/apcupsd_exporter" ]