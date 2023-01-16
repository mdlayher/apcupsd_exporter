package apcupsdexporter

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func testCollector(t *testing.T, collector prometheus.Collector) []byte {
	t.Helper()

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(collector)

	srv := httptest.NewServer(promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	defer srv.Close()

	c := &http.Client{Timeout: 1 * time.Second}
	resp, err := c.Get(srv.URL)
	if err != nil {
		t.Fatalf("failed to HTTP GET data from prometheus: %v", err)
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read server response: %v", err)
	}

	return buf
}
