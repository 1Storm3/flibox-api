package metrics

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	appName   = "kinopoisk_api"
	namespace = "kinopoisk_api__"
	help      = "Total number of requests"
)

type Metrics struct {
	requestCounter prometheus.Counter
}

var metrics *Metrics

func Init(_ context.Context) error {
	metrics = &Metrics{
		requestCounter: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      appName,
				Subsystem: "http",
				Help:      help,
			}),
	}
	return nil
}

func IncRequestCounter() {
	metrics.requestCounter.Inc()
}
