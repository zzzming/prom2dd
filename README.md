# Sending Astra Streaming tenant Prometheus metrics to DataDog

This guide outlines the process for transmitting key gauge metrics from Astra Streaming tenants to DataDog.

The procedure involves scraping Prometheus metrics on a per-tenant basis. To initiate this process, the following two environment variables must be set with the specific tenant scraping URL and JWT:

TARGET_URL: The URL for tenant-specific Prometheus scraping.
PROMETHEUS_JWT_HEADER: The JWT header required for authentication.
For further details on scraping Astra Streaming metrics, refer to the official documentation: [Astra Streaming - Scraping Metrics](https://docs.datastax.com/en/streaming/astra-streaming/operations/astream-scrape-metrics.html).

Additionally, integration with DataDog requires setting up DataDog API and APP keys in your environment. Define these variables as follows:

DD_SITE: Set this to "datadoghq.com".
DD_API_KEY: Your unique DataDog API key.
DD_APP_KEY: Your DataDog application key.
For more information on submitting metrics to DataDog, visit the DataDog documentation: [Submit Metrics to DataDog](https://docs.datadoghq.com/api/latest/metrics/#submit-metrics).

## How to build
```
go build -o prom2dd -tags musl src/main.go
```

## Build image
```
make container
```