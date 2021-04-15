package exporter

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/showwin/speedtest-go/speedtest"
)

const (
	namespace = "speedtest"
)

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last speedtest successful.",
		nil, nil,
	)
	scrapeDurationSeconds = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "scrape_duration_seconds"),
		"Time to preform last speed test",
		nil, nil,
	)
	latency = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "latency_seconds"),
		"Measured latency on last speed test",
		[]string{"user_lat", "user_lon", "user_ip", "user_isp", "server_lat", "server_lon", "server_id", "server_name", "server_country", "distance"},
		nil,
	)
	upload = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "upload_speed_Bps"),
		"Last upload speedtest result",
		[]string{"user_lat", "user_lon", "user_ip", "user_isp", "server_lat", "server_lon", "server_id", "server_name", "server_country", "distance"},
		nil,
	)
	download = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "download_speed_Bps"),
		"Last download speedtest result",
		[]string{"user_lat", "user_lon", "user_ip", "user_isp", "server_lat", "server_lon", "server_id", "server_name", "server_country", "distance"},
		nil,
	)
)

// Exporter runs speedtest and exports them using
// the prometheus metrics package.
type Exporter struct{}

// New returns an initialized Exporter.
func New() (*Exporter, error) {
	return &Exporter{}, nil
}

// Describe describes all the metrics. It implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- scrapeDurationSeconds
	ch <- latency
	ch <- upload
	ch <- download
}

// Collect fetches the stats from Starlink dish and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	ok := e.speedtest(ch)
	d := time.Since(start).Seconds()

	if ok {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 1.0,
		)
		ch <- prometheus.MustNewConstMetric(
			scrapeDurationSeconds, prometheus.GaugeValue, d,
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0.0,
		)
	}
}

func (e *Exporter) speedtest(ch chan<- prometheus.Metric) bool {
	user, err := speedtest.FetchUserInfo()
	if err != nil {
		log.Errorf("could not fetch user information: %s", err.Error())
		return false
	}

	// returns list of servers in distance order
	serverList, err := speedtest.FetchServerList(user)
	if err != nil {
		log.Errorf("could not fetch server list: %s", err.Error())
		return false
	}
	// taking the closes server
	servers := serverList.Servers
	server := servers[0]

	if err := server.PingTest(); err == nil {
		ch <- prometheus.MustNewConstMetric(
			latency, prometheus.GaugeValue, server.Latency.Seconds(),
			user.Lat,
			user.Lon,
			user.IP,
			user.Isp,
			server.Lat,
			server.Lon,
			server.ID,
			server.Name,
			server.Country,
			fmt.Sprintf("%f", server.Distance),
		)
	}

	if err := server.DownloadTest(false); err == nil {
		ch <- prometheus.MustNewConstMetric(
			download, prometheus.GaugeValue, server.DLSpeed*1024*1024,
			user.Lat,
			user.Lon,
			user.IP,
			user.Isp,
			server.Lat,
			server.Lon,
			server.ID,
			server.Name,
			server.Country,
			fmt.Sprintf("%f", server.Distance),
		)
	}

	if err := server.UploadTest(false); err == nil {
		ch <- prometheus.MustNewConstMetric(
			upload, prometheus.GaugeValue, server.ULSpeed*1024*1024,
			user.Lat,
			user.Lon,
			user.IP,
			user.Isp,
			server.Lat,
			server.Lon,
			server.ID,
			server.Name,
			server.Country,
			fmt.Sprintf("%f", server.Distance),
		)
	}

	return true
}
