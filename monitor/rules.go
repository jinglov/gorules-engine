package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

var (
	RuleCodeVec  = NewRuleVec("go_rules", "engine", "rule_code")
	RuleCidVec   = NewRuleVec("go_rules", "engine", "rule_cid")
	ruleLabelMap *sync.Map
)

type ruleVec struct {
	vec *prometheus.CounterVec
}

type RuleLabels struct {
	Key             string
	prometheusLabel prometheus.Labels
	sync.Once
}

func init() {
	ruleLabelMap = new(sync.Map)
}

func GetRuleLabels(key string) *RuleLabels {
	if val, ok := ruleLabelMap.Load(key); ok {
		return val.(*RuleLabels)
	}
	obj := &RuleLabels{
		Key: key,
	}
	val, _ := ruleLabelMap.LoadOrStore(key, obj)
	return val.(*RuleLabels)
}

func (l *RuleLabels) toPrometheusLable() prometheus.Labels {
	l.Do(func() {
		l.prometheusLabel = prometheus.Labels{
			"key": l.Key,
		}
	})
	return l.prometheusLabel
}

func NewRuleVec(namespace, subsystem, name string) *ruleVec {
	return &ruleVec{
		vec: NewCounterVec(namespace, subsystem, name, "go rules engine rule counter by handlers", []string{"key"}),
	}
}

func (kv *ruleVec) Inc(labels *RuleLabels) {
	kv.vec.With(labels.toPrometheusLable()).Inc()
}
