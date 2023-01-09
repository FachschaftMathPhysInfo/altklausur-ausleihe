package utils

import (
	"github.com/FachschaftMathPhysInfo/altklausur-ausleihe/server/graph/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"
)

var (
	// ExamsMarkedMetric is a prometheus metric for the amount of requested exams for watermarking
	ExamsMarkedMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name: "altklausur_ausleihe_exams_marked",
		Help: "The total number of requested exams",
	})

	// TotalExamsMetric is a prometheus metric for the total amount of exams in the database
	TotalExamsMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "altklausur_ausleihe_exams_total",
		Help: "The total number of exams",
	})
)

// GetTotalExamsMetric updates the TotalExamsMetric with the current value from the database
func GetTotalExamsMetric(database *gorm.DB) {
	var count int64
	database.Model(&model.Exam{}).Count(&count)
	TotalExamsMetric.Set(float64(count))
}
