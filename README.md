<p align=center>
  <img alt=logo src="https://github.com/danopstech/speedtest_exporter/raw/main/.docs/assets/logo.jpg" height=150 />
  <h3 align=center>Speedtest Prometheus Exporter</h3>
</p>

---
A [Speedtest](https://www.speedtest.net) exporter for Prometheus.

[![goreleaser](https://github.com/danopstech/speedtest_exporter/actions/workflows/release.yaml/badge.svg)](https://github.com/danopstech/speedtest_exporter/actions/workflows/release.yaml)
[![License](https://img.shields.io/github/license/danopstech/speedtest_exporter)](/LICENSE)
[![Release](https://img.shields.io/github/release/danopstech/speedtest_exporter.svg)](https://github.com/danopstech/speedtest_exporter/releases/latest)
[![Docker](https://img.shields.io/docker/pulls/danopstech/speedtest_exporter)](https://hub.docker.com/r/danopstech/speedtest_exporter)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/danopstech/speedtest_exporter)

## Simple Usage:

### Flags

`speedtest_exporter` is configured by optional command line flags

```bash
$ ./speedtest_exporter --help
Usage of speedtest_exporter
  -port string
        Listening port for prometheus endpoint (default 9090)

```

### Binaries

For pre-built binaries please take a look at the [releases](https://github.com/danopstech/speedtest_exporter/releases).

```bash
./speedtest_exporter [flags]
```

### Docker

Docker Images can be found at [GitHub Container Registry](https://github.com/orgs/danopstech/packages/container/package/speedtest_exporter) & [Dockerhub](https://hub.docker.com/r/danopstech/speedtest_exporter).

Example:
```bash
docker pull ghcr.io/danopstech/speedtest_exporter:latest

docker run \
  -p 9090:9090 \
  ghcr.io/danopstech/speedtest_exporter:latest [flags]
```

### Setup Prometheus to scrape `speedtest_exporter`

Configure [Prometheus](https://prometheus.io/) to scrape metrics from localhost:9090/metrics

This exporter locks (one concurrent scrape at a time) as it conducts the speedtest when scraped, **remember set scrape interval, and scrap timeout** accordingly as per example.

```yaml
...
scrape_configs
    - job_name: speedtest
      scrape_interval: 60m
      scrape_timeout:  40s
      static_configs:
        - targets: ['localhost:9090']
...
```

## Exported Metrics:

```
# HELP speedtest_download_speed_Bps Last download speedtest result
# TYPE speedtest_download_speed_Bps gauge
# HELP speedtest_latency_seconds Measured latency on last speed test
# TYPE speedtest_latency_seconds gauge
# HELP speedtest_scrape_duration_seconds Time to preform last speed test
# TYPE speedtest_scrape_duration_seconds gauge
# HELP speedtest_up Was the last speedtest successful.
# TYPE speedtest_up gauge
# HELP speedtest_upload_speed_Bps Last upload speedtest result
# TYPE speedtest_upload_speed_Bps gauge
```
