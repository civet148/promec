package promec

import (
	"fmt"
	"github.com/civet148/log"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

type Metrics interface {
	Key() string
	Metric() prometheus.Metric
}

type ConstMetric struct {
	name   string
	labels []string
	values []string
	metric prometheus.Metric
}

func newConstMetric(name string, metric prometheus.Metric, labelNames, labelValues []string) *ConstMetric {
	if len(labelNames) == 0 || len(labelValues) == 0 {
		log.Panic("nil labels or values")
	}
	if len(labelNames) != len(labelValues) {
		log.Panic("label and value size not equal")
	}
	return &ConstMetric{
		name:   name,
		metric: metric,
		labels: labelNames,
		values: labelValues,
	}
}

func (m *ConstMetric) Key() string {
	var strKey string
	var lvs []string
	for i, v := range m.labels {
		lvs = append(lvs, fmt.Sprintf("%s=\"%s\"", v, m.values[i]))
	}
	strKey = fmt.Sprintf("%s{%s}", m.name, strings.Join(lvs, ","))
	return strKey
}

func (m *ConstMetric) Metric() prometheus.Metric {
	return m.metric
}
