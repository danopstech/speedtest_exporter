package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/danopstech/speedtest_exporter/internal/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	metricsPath = "/metrics"
)

func main() {
	port := flag.String("port", "9090", "listening port to expose metrics on")
	serverID := flag.Int("server_id", -1, "Speedtest.net server ID to run test against, -1 will pick the closest server to your location")
	serverFallback := flag.Bool("server_fallback", false, "If the server_id given is not available, should we fallback to closest available server")
	requestTimeout := flag.Int("timeout", 60, "request timeout for the execution of the speedtest")
	flag.Parse()

	exporter, err := exporter.New(*serverID, *serverFallback)
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
             <p>Metrics page will take approx 40 seconds to load and show results, as the exporter carries out a speedtest when scraped.</p>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
             <p><a href='/health'>Health</a></p>
             </body>
             </html>`))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		client := http.Client{
			Timeout: 3 * time.Second,
		}
		_, err := client.Get("https://clients3.google.com/generate_204")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "No Internet Connection")
		} else {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, "OK")
		}
	})

	http.Handle(metricsPath, promhttp.HandlerFor(r, promhttp.HandlerOpts{
		MaxRequestsInFlight: 1,
		Timeout:             time.Duration(*requestTimeout) * time.Second,
	}))

	log.Info("starting listener on port: " + *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
