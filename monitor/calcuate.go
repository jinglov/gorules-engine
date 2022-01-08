package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

const (
	HFUNC     = "hfunc"
	LFUNC     = "lfunc"
	CACHEFUNC = "cachefunc"
)

var actCalLabels = new(sync.Map)

type CalcuateVec struct {
	vec *prometheus.CounterVec
}

type CalcuateLabels struct {
	Key             string
	FunType         string
	prometheusLabel prometheus.Labels
	sync.Once
}

type ActionCalcuateLabels struct {
	HFunc     *CalcuateLabels
	LFunc     *CalcuateLabels
	CacheFunc *CalcuateLabels
}

func GetActionCalcuateLabels(key string) *ActionCalcuateLabels {
	if v, ok := actCalLabels.Load(key); ok {
		return v.(*ActionCalcuateLabels)
	}

	labels := NewActionCalcuateLabels(key)
	actCalLabels.Store(key, labels)
	return labels
}

func NewActionCalcuateLabels(key string) *ActionCalcuateLabels {
	return &ActionCalcuateLabels{
		HFunc:     &CalcuateLabels{Key: key, FunType: HFUNC},
		LFunc:     &CalcuateLabels{Key: key, FunType: LFUNC},
		CacheFunc: &CalcuateLabels{Key: key, FunType: CACHEFUNC},
	}
}

func (l *CalcuateLabels) toPrometheusLable() prometheus.Labels {
	l.Do(func() {
		l.prometheusLabel = prometheus.Labels{
			"key":      l.Key,
			"fun_type": l.FunType,
		}
	})

	return l.prometheusLabel
}

var (
	SaverCalVec      = NewCalcuateVec("go_rules", "engine", "calcuate_saver")
	RulerCalVec      = NewCalcuateVec("go_rules", "engine", "calcuate_ruler")
	OutputCalVec     = NewCalcuateVec("go_rules", "engine", "calcuate_output")
	OutputCodeCalVec = NewCalcuateVec("go_rules", "engine", "calcuate_output_code")
)

func NewCalcuateVec(namespace, subsystem, name string) *CalcuateVec {
	return &CalcuateVec{
		vec: NewCounterVec(namespace, subsystem, name, "go rules engine func calcuate num", []string{"key", "fun_type"}),
	}
}

func (kv *CalcuateVec) Add(labels *CalcuateLabels, n int) {
	kv.vec.With(labels.toPrometheusLable()).Add(float64(n))
}
