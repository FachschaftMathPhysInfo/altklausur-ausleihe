package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// WatermarkingTimeHistogram is a prometheus metric for the total amount of exams in the database
	WatermarkingTimeHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "altklausur_ausleihe_watermarking_histogram",
		Help: "histogram for the time it takes an exam to get marked",
	})
)
