package monitor

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	OutputCidVec  = NewOutputRuleVec("go_rules", "engine", "output_cid")
	OutputCodeVec = NewOutputRuleVec("go_rules", "engine", "output_code")
	OutputTagVec  = NewOutputRuleVec("go_rules", "engine", "output_tag")
	OutputKeyVec  = NewOutputRuleVec("go_rules", "engine", "output_key")
	outputRuleMap *sync.Map
)

type outputRuleVec struct {
	vec *prometheus.CounterVec
}

func init() {
	outputRuleMap = new(sync.Map)
}

func GetOutputRuleLabel(key string) *OutputRuleLabels {
	if val, ok := outputRuleMap.Load(key); ok {
		return val.(*OutputRuleLabels)
	}
	obj := &OutputRuleLabels{
		Key: key,
	}
	val, _ := outputRuleMap.LoadOrStore(key, obj)
	return val.(*OutputRuleLabels)
}

type OutputRuleLabels struct {
	Key             string
	prometheusLabel prometheus.Labels
	sync.Once
}

func (l *OutputRuleLabels) toPrometheusLable() prometheus.Labels {
	l.Do(func() {
		l.prometheusLabel = prometheus.Labels{
			"key": l.Key,
		}
	})
	return l.prometheusLabel
}

func NewOutputRuleVec(namespace, subsystem, name string) *outputRuleVec {
	return &outputRuleVec{
		vec: NewCounterVec(namespace, subsystem, name, "go rules engine output counter by handlers", []string{"key"}),
	}
}

func (kv *outputRuleVec) Inc(labels *OutputRuleLabels) {
	kv.vec.With(labels.toPrometheusLable()).Inc()
}
