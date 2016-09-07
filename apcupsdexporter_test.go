package apcupsdexporter

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func testCollector(t *testing.T, collector prometheus.Collector) []byte {
	if err := prometheus.Register(collector); err != nil {
		t.Fatalf("failed to register Prometheus collector: %v", err)
	}
	defer prometheus.Unregister(collector)

	promServer := httptest.NewServer(prometheus.Handler())
	defer promServer.Close()

	resp, err := http.Get(promServer.URL)
	if err != nil {
		t.Fatalf("failed to GET data from prometheus: %v", err)
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read server response: %v", err)
	}

	return buf
}
