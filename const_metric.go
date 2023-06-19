package promec

import (
	"fmt"
	"github.com/civet148/log"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"sync"
)

type Metrics interface {
	Key() string               //metric key
	Get() float64              //get metric value
	Set(value float64)         //set metric value
	Metric() prometheus.Metric //get metric object
}

type ConstMetric struct {
	name        string
	labelNames  []string
	labelValues []string
	value       float64
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
		value:       value,
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

func (m *ConstMetric) Get() float64 {
	m.locker.RLock()
	defer m.locker.RUnlock()
	return m.value
}

func (m *ConstMetric) Set(value float64) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.value = value
	m.renewMetricNoLock()
}

func (m *ConstMetric) Add(value float64) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.value += value
	m.renewMetricNoLock()
}

func (m *ConstMetric) Inc() {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.value += 1
	m.renewMetricNoLock()
}

func (m *ConstMetric) Dec() {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.value -= 1
	m.renewMetricNoLock()
}

func (m *ConstMetric) renewMetricNoLock() {
	m.metric = prometheus.MustNewConstMetric(m.desc, m.vt, m.value, m.labelValues...)
}
