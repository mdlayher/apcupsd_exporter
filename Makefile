VERSION = v0.2.0

binary:
	go build cmd/apcupsd_exporter/main.go -o apcupsd_exporter

buildah:
	buildah bud -t apcupsd_exporter:${VERSION}