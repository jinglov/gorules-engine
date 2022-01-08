package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

type FuncVec struct {
	vec *prometheus.CounterVec
}

type FuncLabel struct {
	Name            string
	prometheusLabel prometheus.Labels
	sync.Once
}

func (l *FuncLabel) toPrometheusLable() prometheus.Labels {
	l.Do(func() {
		l.prometheusLabel = prometheus.Labels{
			"name": l.Name,
		}
	})
	return l.prometheusLabel
}

func NewFuncVec(namespace, subsystem, name string) *FuncVec {
	return &FuncVec{
		vec: NewCounterVec(namespace, subsystem, name, "go rules engine counter by handlers", []string{"name"}),
	}
}

func (kv *FuncVec) Inc(labels *FuncLabel) {
	kv.vec.With(labels.toPrometheusLable()).Inc()
}
