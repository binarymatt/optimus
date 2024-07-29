package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	records_processed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "records_processed",
	}, []string{
		"type",
		"name",
	})
)

func RecordProcessedRecord(processor_type, name string) {
	records_processed.WithLabelValues(processor_type, name).Inc()
}
