package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type saverVec struct {
	vec *prometheus.CounterVec
}

type SaverLabels struct {
	DataSource      string
	ResDataKey      string
	Status          string
	prometheusLabel prometheus.Labels
	sync.Once
}

var (
	SaverVec = NewSaverVec("go_rules", "engine", "saver")
)

func (l *SaverLabels) toPrometheusLable() prometheus.Labels {
	l.Do(func() {
		l.prometheusLabel = prometheus.Labels{
			"datasource": l.DataSource,
			"key":        l.ResDataKey,
			"status":     l.Status,
		}
	})

	return l.prometheusLabel
}

func NewSaverVec(namespace, subsystem, name string) *saverVec {
	return &saverVec{
		vec: NewCounterVec(namespace, subsystem, name, "go rules saver counter by handlers", []string{"datasource", "key", "status"}),
	}
}

func (kv *saverVec) Inc(labels *SaverLabels) {
	kv.vec.With(labels.toPrometheusLable()).Inc()
}
