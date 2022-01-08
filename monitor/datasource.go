package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var (
	DSVec         = NewDSVec("go_rules", "engine", "datasource")
	datasourceMap *sync.Map
)

func init() {
	datasourceMap = new(sync.Map)
}

type dsVec struct {
	vec *prometheus.CounterVec
}

type DSLabel struct {
	DataSource      string
	Step            string
	prometheusLabel prometheus.Labels
	sync.Once
}

func GetDSLabel(datasource, step string) *DSLabel {
	if val, ok := datasourceMap.Load(datasource + step); ok {
		return val.(*DSLabel)
	}
	obj := &DSLabel{
		DataSource: datasource,
		Step:       step,
	}
	val, _ := datasourceMap.LoadOrStore(datasource+step, obj)
	return val.(*DSLabel)
}

func (l *DSLabel) toPrometheusLable() prometheus.Labels {
	l.Do(func() {
		l.prometheusLabel = prometheus.Labels{
			"datasource": l.DataSource,
			"step":       l.Step,
		}
	})
	return l.prometheusLabel
}

func NewDSVec(namespace, subsystem, name string) *dsVec {
	return &dsVec{
		vec: NewCounterVec(namespace, subsystem, name, "go rules engine datasource counter by handlers", []string{"datasource", "step"}),
	}
}

func (kv *dsVec) Inc(labels *DSLabel) {
	kv.vec.With(labels.toPrometheusLable()).Inc()
}
