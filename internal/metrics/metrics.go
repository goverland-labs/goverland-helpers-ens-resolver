package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "inbox",
			Name:      "sdk_requests",
			Help:      "Time taken to process requests",
			Buckets:   []float64{.005, .01, .025, .05, .075, .1, .15, .2, .25, .5, 1, 2.5, 5, 10, 15, 30},
		},
		[]string{"client", "method", "error"},
	)
)

func CollectRequestsMetric(client, method string, err error, start time.Time) {
	RequestsHistogram.
		WithLabelValues(client, method, errLabelValue(err)).
		Observe(time.Since(start).Seconds())
}

// ErrLabelValue returns string representation of error label value
func errLabelValue(err error) string {
	if err != nil {
		return "true"
	}

	return "false"
}
