package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"

	"github.com/danopstech/speedtest_exporter/internal/exporter"
)

const (
	metricsPath = "/metrics"
)

func main() {
	port := flag.String("port", "9090", "listening port to expose metrics on")
	flag.Parse()

	exporter, err := exporter.New()
	if err != nil {
		panic(err)
	}

	r := prometheus.NewRegistry()
	r.MustRegister(exporter)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
             <head><title>Speedtest Exporter</title></head>
             <body>
             <h1>Speedtest Exporter</h1>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "OK")
	})

	http.Handle(metricsPath, promhttp.HandlerFor(r, promhttp.HandlerOpts{
		MaxRequestsInFlight: 1,
		Timeout:             40 * time.Second,
	}))
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
