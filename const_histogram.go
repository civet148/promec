package promec

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"strings"
	"sync"
)

type ConstHistogram struct {
	name        string
	count       uint64
	sum         float64
	labelNames  []string
	labelValues []string
	locker      sync.RWMutex
	buckets     map[float64]uint64
	desc        *prometheus.Desc
	metric      prometheus.Metric
}

// NewConstHistogram creates a new ConstHistogram
// count: total count of sampling
// sum: total sum of sampling
func newConstHistogram(name string, desc *prometheus.Desc, count uint64, sum float64, buckets map[float64]uint64, labelNames, labelValues []string) *ConstHistogram {
	metric := prometheus.MustNewConstHistogram(desc, count, sum, buckets, labelValues...)
	return &ConstHistogram{
		name:        name,
		desc:        desc,
		sum:         sum,
		count:       count,
		metric:      metric,
		buckets:     buckets,
		labelNames:  labelNames,
		labelValues: labelValues,
	}
}

func (m *ConstHistogram) Key() string {
	var strKey string
	var lvs []string
	for i, v := range m.labelNames {
		lvs = append(lvs, fmt.Sprintf("%s=\"%s\"", v, m.labelValues[i]))
	}
	strKey = fmt.Sprintf("%s{%s}", m.name, strings.Join(lvs, ","))
	return strKey
}

func (m *ConstHistogram) Metric() prometheus.Metric {
	m.locker.Lock()
	defer m.locker.Unlock()
	return m.metric
}

func (m *ConstHistogram) Set(count uint64, sum float64, buckets map[float64]uint64) {
	m.locker.Lock()
	defer m.locker.Unlock()
	m.sum = sum
	m.count = count
	m.buckets = buckets
	m.renewMetricNoLock()
}

func (m *ConstHistogram) renewMetricNoLock() {
	m.metric = prometheus.MustNewConstHistogram(m.desc, m.count, m.sum, m.buckets, m.labelValues...)
}
