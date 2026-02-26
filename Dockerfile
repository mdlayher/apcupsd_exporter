FROM golang:1.26-alpine AS build

WORKDIR /app

COPY . .

RUN go build ./cmd/apcupsd_exporter

FROM gcr.io/distroless/static:nonroot AS runtime

COPY --from=build /app/apcupsd_exporter /apcupsd_exporter

WORKDIR /

CMD ["/apcupsd_exporter"]

