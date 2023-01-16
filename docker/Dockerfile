FROM golang:latest

COPY ./src /src
WORKDIR /src

RUN go install && \
    CGO_ENABLED=0 go build -o apcups-exporter cmd/apcupsd_exporter/main.go && \
    chmod +x apcups-exporter

FROM alpine:latest

COPY --from=0 /src/apcups-exporter /usr/sbin/apcups-exporter

RUN adduser -s /bin/sh -H -D apcexp

USER apcexp

ENTRYPOINT ["apcups-exporter"]
CMD ["--help"]