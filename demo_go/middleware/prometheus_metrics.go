package middleware

import (
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/wojciech12/order-manager/service"
)

// might go to package middleware
type PromMetricCollector struct {
	aboutMe       *prometheus.HistogramVec
	aboutDatabase *prometheus.HistogramVec
	aboutAudit    *prometheus.HistogramVec
}

func (m *PromMetricCollector) ObserveMyself(path string, method string, statusCode int, duration float64) {
	sc := strconv.Itoa(statusCode)
	m.aboutMe.WithLabelValues(path, method, sc).Observe(duration)
}

func (m *PromMetricCollector) ObserveDatabase(queryId string, statusCode int, sqlState string, duration float64) {
	sc := strconv.Itoa(statusCode)
	m.aboutDatabase.WithLabelValues(queryId, sc, sqlState).Observe(duration)
}

func (m *PromMetricCollector) ObserveAudit(path string, method string, statusCode int, duration float64) {
	sc := strconv.Itoa(statusCode)
	m.aboutAudit.WithLabelValues(path, method, sc).Observe(duration)
}

func NewMetricCollector(serviceName string) service.MetricCollector {
	prefix := strings.Replace(serviceName, "-", "_", -1)

	m := PromMetricCollector{}

	m.aboutMe = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: prefix + "_duration_seconds",
			Help: serviceName + " latency request distribution",
		},
		[]string{"path", "method", "status_code"},
	)

	m.aboutDatabase = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: prefix + "_database_duration_seconds",
			Help: "database latency request distribution",
		},
		[]string{"queryId", "statusCode", "sqlState"},
	)

	m.aboutAudit = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: prefix + "_audit_duration_seconds",
			Help: serviceName + " latency request distribution",
		},
		[]string{"path", "method", "status_code"},
	)
	m.registerHandlers()
	return &m
}

func (m *PromMetricCollector) registerHandlers() {
	prometheus.MustRegister(m.aboutMe)
	prometheus.MustRegister(m.aboutDatabase)
	prometheus.MustRegister(m.aboutAudit)
}
