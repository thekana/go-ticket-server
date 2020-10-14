package prometheus

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpAPICallDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: prometheusNamespace,
		Subsystem: "http_api",
		Name:      "call_duration_seconds",
		Help:      "Duration of HTTP call in seconds",
		Buckets:   prometheus.ExponentialBuckets(0.0001, 1.5, 36),
	},
		[]string{"code", "path"},
	)
)

func RecordHTTPAPICallDurationMetrics(duration time.Duration, code int, path string) {
	if config == nil || !config.Enable {
		return
	}

	httpAPICallDurationHistogram.WithLabelValues(
		strconv.Itoa(code),
		path,
	).Observe(duration.Seconds())
}
