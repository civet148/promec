package promec

import (
	"fmt"
	"github.com/civet148/log"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"sync"
)

type Metrics interface {
	Key() string
	Set(value float64)
	Metric() prometheus.Metric
}

type ConstMetric struct {
	name        string
	labelNames  []string
	labelValues []string
	vt          prometheus.ValueType
	desc        *prometheus.Desc
	locker      sync.RWMutex
	metric      prometheus.Metric
}

func newConstMetric(name string, vt prometheus.ValueType, desc *prometheus.Desc, value float64, labelNames, labelValues []string) *ConstMetric {
	if len(labelNames) != len(labelValues) {
		log.Panic("label and value size not equal")
	}
	var metric prometheus.Metric
	metric = prometheus.MustNewConstMetric(desc, vt, value, labelValues...)
	return &ConstMetric{
		name:        name,
		desc:        desc,
		vt:          vt,
		metric:      metric,
		labelNames:  labelNames,
		labelValues: labelValues,
	}
}

func (m *ConstMetric) Key() string {
	var strKey string
	var lvs []string
	for i, v := range m.labelNames {
		lvs = append(lvs, fmt.Sprintf("%s=\"%s\"", v, m.labelValues[i]))
	}
	strKey = fmt.Sprintf("%s{%s}", m.name, strings.Join(lvs, ","))
	return strKey
}

func (m *ConstMetric) Metric() prometheus.Metric {
	m.locker.RLock()
	defer m.locker.RUnlock()
	return m.metric
}

func (m *ConstMetric) Set(value float64) {
	m.locker.Lock()
	defer m.locker.Unlock()
	metric := prometheus.MustNewConstMetric(m.desc, m.vt, value, m.labelValues...)
	m.metric = metric
}
