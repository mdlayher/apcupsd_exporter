apcupsd_exporter [![Build Status](https://travis-ci.org/mdlayher/apcupsd_exporter.svg?branch=master)](https://travis-ci.org/mdlayher/apcupsd_exporter) [![GoDoc](http://godoc.org/github.com/mdlayher/apcupsd_exporter?status.svg)](http://godoc.org/github.com/mdlayher/apcupsd_exporter)
================

Command `apcupsd_exporter` provides a Prometheus exporter for the
[apcupsd](http://www.apcupsd.org/) Network Information Server (NIS).
MIT Licensed.

Usage
-----

Available flags for `apcupsd_exporter` include:

```
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

Sample
------

Here is a screenshot of an example grafana dashboard using metrics from `apcupsd_exporter`.
You can find the dashboard JSON in the `grafana` folder.
![sample](https://raw.githubusercontent.com/jangrewe/apcupsd_exporter/master/grafana/screenshot.png)

