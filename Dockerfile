FROM gcr.io/distroless/static
ENTRYPOINT ["/speedtest_exporter"]
COPY speedtest_exporter /
