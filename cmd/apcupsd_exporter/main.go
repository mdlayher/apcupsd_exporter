// Command apcupsd_exporter provides a Prometheus exporter for apcupsd.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/mdlayher/apcupsd"
	apcupsdexporter "github.com/mdlayher/apcupsd_exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	apcupsdAddr    = flag.String("apcupsd.addr", ":3551", "address of apcupsd Network Information Server (NIS)")
	apcupsdNetwork = flag.String("apcupsd.network", "tcp", `network of apcupsd Network Information Server (NIS): typically "tcp", "tcp4", or "tcp6"`)

	telemetryAddr = flag.String("telemetry.addr", ":9162", "address for apcupsd exporter")
	metricsPath   = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics")
)

func main() {
	flag.Parse()

	if telemetryAddrArg := os.Getenv("TELEMETRY_PATH"); telemetryAddrArg != "" {
		telemetryAddr = &telemetryAddrArg
	} 
	if metricsPathArg := os.Getenv("METRICS_PATH"); metricsPathArg != "" {
		metricsPath = &metricsPathArg
	}
	if apcupsdAddrArg := os.Getenv("APCUPSD_ADDRESS"); apcupsdAddrArg != "" {
		apcupsdAddr = &apcupsdAddrArg
	}
	if apcupsdNetworkArg := os.Getenv("APCUPSD_NETWORK"); apcupsdNetworkArg != "" {
		apcupsdNetwork = &apcupsdNetworkArg
	}
	
	if *apcupsdAddr == "" {
		log.Fatal("address of apcupsd Network Information Server (NIS) must be specified with '-apcupsd.addr' flag")
	}

	fn := newClient(*apcupsdNetwork, *apcupsdAddr)

	prometheus.MustRegister(apcupsdexporter.New(fn))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Printf("starting apcupsd exporter on %q for server %s://%s",
		*telemetryAddr, *apcupsdNetwork, *apcupsdAddr)

	if err := http.ListenAndServe(*telemetryAddr, nil); err != nil {
		log.Fatalf("cannot start apcupsd exporter: %s", err)
	}
}

func newClient(network, addr string) apcupsdexporter.ClientFunc {
	return func(ctx context.Context) (*apcupsd.Client, error) {
		return apcupsd.DialContext(ctx, network, addr)
	}
}
