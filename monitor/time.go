package monitor

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	TimeVec = NewTimeVec("go_rules", "engine", "time",
		[]float64{
			5, 10, 20, 50, 100, 500, 1000, 3000})
	timeMap *sync.Map
)

type timeVec struct {
	vec *prometheus.HistogramVec
}

func init() {
	timeMap = new(sync.Map)
}

func GetTimeLabel(dataSource string) *TimeLabels {
	if val, ok := timeMap.Load(dataSource); ok {
		return val.(*TimeLabels)
	}
	obj := &TimeLabels{
		Datasource: dataSource,
	}
	val, _ := timeMap.LoadOrStore(dataSource, obj)
	return val.(*TimeLabels)
}

type TimeLabels struct {
	Datasource      string
	prometheusLabel prometheus.Labels
	sync.Once
}

func (l *TimeLabels) toPrometheusLable() prometheus.Labels {
	l.Do(func() {
		l.prometheusLabel = prometheus.Labels{
			"datasource": l.Datasource,
		}
	})
	return l.prometheusLabel
}

func NewTimeVec(namespace, subsystem, name string, bucktes []float64) *timeVec {
	return &timeVec{
		vec: NewHistogramVec(namespace, subsystem, name, "go rules engine entire process time", []string{"datasource"}, bucktes),
	}
}

func (kv *timeVec) Observe(labels *TimeLabels, elapsed float64) {
	kv.vec.With(labels.toPrometheusLable()).Observe(elapsed)
}
