package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
	"github.com/caarlos0/env/v9"

	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"
)

// DD_SITE="datadoghq.com" DD_API_KEY="<DD_API_KEY>" DD_APP_KEY="<DD_APP_KEY>"
const (
	defaultGauges = "pulsar_msg_backlog"
)

// Config is the configuration for the exporter
type Config struct {
	PrometheusScrapeURL    string `env:"TARGET_URL"`
	ScrapeInterval         int    `env:"SCRAPE_INTERVAL" envDefault:"60"`
	PrometheusBearerHeader string `env:"PROMETHEUS_JWT_HEADER"`
	Metrics                string `env:"METRICS" envDefault:"pulsar_msg_backlog"`
}

// DD_SITE="datadoghq.com" DD_API_KEY="<DD_API_KEY>" DD_APP_KEY="<DD_APP_KEY>"

func main() {
	// Create a channel to receive OS signals.
	sigs := make(chan os.Signal, 1)

	// Notify the channel for SIGINT (Ctrl+C) and SIGTERM signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Config deserialize error %v", err)
	}

	// Start an infinite loop.
	go func(c Config) {
		log.Println("Starting...")
		ticker := time.NewTicker(time.Duration(c.ScrapeInterval) * time.Second)
		defer ticker.Stop()

		metrics := strings.Split(c.Metrics, ",")
		for {
			select {
			case <-ticker.C:
				err := scrapePrometheus(c.PrometheusScrapeURL, c.PrometheusBearerHeader, metrics)
				if err != nil {
					log.Printf("Error scraping Prometheus %s: %v", c.PrometheusBearerHeader, err)
				}
			}
		}
	}(cfg)

	<-sigs
}

func scrapePrometheus(targetURL, token string, metrics []string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Add the Bearer token for authorization
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error fetching metrics: %v", err)
	}
	defer resp.Body.Close()

	var parser expfmt.TextParser
	metricFamilies, err := parser.TextToMetricFamilies(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing metrics: %v", err)
	}

	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewMetricsApi(apiClient)

	// Iterate through metrics and print them
	for metricName, metricFamily := range metricFamilies {
		if !contains(metrics, metricName) {
			continue
		}
		for _, metric := range metricFamily.Metric {
			labels := make(model.LabelSet)
			for _, label := range metric.Label {
				labels[model.LabelName(*label.Name)] = model.LabelValue(*label.Value)
			}

			// For demonstration, just print the metric names, labels, and counters.
			// Adjust this section as needed.
			var value float64
			if metric.Gauge != nil {
				value = metric.Gauge.GetValue()
			} else if metric.Untyped != nil {
				value = metric.Untyped.GetValue()
			} else {
				continue
			}

			// fmt.Printf("Metric: %s | Labels: %v | Value: %f\n", metricName, labels, value)
			var tags []string
			for _, label := range metric.Label {
				tags = append(tags, fmt.Sprintf("%s:%s", *label.Name, *label.Value))
			}

			body := datadogV2.MetricPayload{
				Series: []datadogV2.MetricSeries{
					{
						Metric: metricName,
						Type:   datadogV2.METRICINTAKETYPE_GAUGE.Ptr(),
						Points: []datadogV2.MetricPoint{
							{
								Timestamp: metric.TimestampMs,
								Value:     datadog.PtrFloat64(value),
							},
						},
						Tags: tags,
					},
				},
			}
			_, _, err := api.SubmitMetrics(ctx, body, *datadogV2.NewSubmitMetricsOptionalParameters())

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error when calling `MetricsApi.SubmitMetrics`: %v\n", err)
				return err
			}
		}
	}

	return nil
}

func contains(slice []string, target string) bool {
	for _, value := range slice {
		if value == target {
			return true
		}
	}
	return false
}
