package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

var ChannelVec = NewChannelVec("go_rules", "engine", "ch_len")

type channelVec struct {
	vec *prometheus.GaugeVec
}

func NewChannelVec(namespace, subsystem, name string) *channelVec {
	return &channelVec{
		vec: NewGaugeVec(namespace, subsystem, name, "go rules engine channel length", []string{"name"}),
	}
}

func (cv *channelVec) Inc(name string) {
	cv.vec.With(prometheus.Labels{"name": name}).Inc()
}

func (cv *channelVec) Dec(name string) {
	cv.vec.With(prometheus.Labels{"name": name}).Dec()
}
