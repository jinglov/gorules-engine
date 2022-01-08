package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var (
	HitVec = NewHitVec("go_rules", "engine", "hit",
		[]float64{5, 10, 50, 100, 500, 1000, 5000, 10000})
	hitMap *sync.Map
)

type hitVec struct {
	vec *prometheus.HistogramVec
}

func init() {
	hitMap = new(sync.Map)
}

func GetHitLabel(dataSource, reason string) *HitLabels {
	if val, ok := hitMap.Load(reason + dataSource); ok {
		return val.(*HitLabels)
	}
	obj := &HitLabels{
		ReasonCode: reason,
		Datasource: dataSource,
	}
	val, _ := hitMap.LoadOrStore(reason+dataSource, obj)
	return val.(*HitLabels)
}

type HitLabels struct {
	ReasonCode      string
	Datasource      string
	prometheusLabel prometheus.Labels
	sync.Once
}

func (l *HitLabels) toPrometheusLable() prometheus.Labels {
	l.Do(func() {
		l.prometheusLabel = prometheus.Labels{
			"reason_code": l.ReasonCode,
			"datasource":  l.Datasource,
		}
	})
	return l.prometheusLabel
}

func NewHitVec(namespace, subsystem, name string, bucktes []float64) *hitVec {
	return &hitVec{
		vec: NewHistogramVec(namespace, subsystem, name, "go rules engine output time by handlers", []string{"reason_code", "datasource"}, bucktes),
	}
}

func (kv *hitVec) Observe(labels *HitLabels, elapsed float64) {
	kv.vec.With(labels.toPrometheusLable()).Observe(elapsed)
}
