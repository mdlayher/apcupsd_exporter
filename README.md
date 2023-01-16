# apcupsd_exporter [![Linux Test Status](https://github.com/mdlayher/apcupsd_exporter/workflows/Linux%20Test/badge.svg)](https://github.com/mdlayher/apcupsd_exporter/actions)  [![GoDoc](http://godoc.org/github.com/mdlayher/apcupsd_exporter?status.svg)](http://godoc.org/github.com/mdlayher/apcupsd_exporter)


Command `apcupsd_exporter` provides a Prometheus exporter for the
[apcupsd](http://www.apcupsd.org/) Network Information Server (NIS). MIT
Licensed.

## Usage

Available flags for `apcupsd_exporter` include:

```bash
$ ./apcupsd_exporter -h
Usage of ./apcupsd_exporter:
  -apcupsd.addr string
        address of apcupsd Network Information Server (NIS) (default ":3551")
  -apcupsd.network string
        network of apcupsd Network Information Server (NIS): typically "tcp", "tcp4", or "tcp6" (default "tcp")
  -telemetry.addr string
        address for apcupsd exporter (default ":9162")
  -telemetry.path string
        URL path for surfacing collected metrics (default "/metrics")
```

## Docker Usage

Per default the docker container runs `--help` start the apcupsd_exporter like this (e.g.)

```bash
docker run \
      -it --rm -p 9162:9162 \
      --add-host host.docker.internal:host-gateway \
      ghcr.io/mdlayher/apcupsd-exporter:latest -apcupsd.addr host.docker.internal:3551 
```
