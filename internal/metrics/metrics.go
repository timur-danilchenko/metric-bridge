package metrics

import "github.com/prometheus/client_golang/prometheus"

var ProcessedMessages = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "processed_total",
	Help: "Total number of processed metrics",
})

var ProcessingErrors = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "errors_total",
	Help: "Total number of errors",
})

var ProcessingDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "processing_duration_seconds",
	Help:    "Metric processing time",
	Buckets: prometheus.DefBuckets,
})

func Register() {
	prometheus.MustRegister(ProcessedMessages)
	prometheus.MustRegister(ProcessingErrors)
	prometheus.MustRegister(ProcessingDuration)
}
