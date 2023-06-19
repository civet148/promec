package promec

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	TagNameLabel   = "label"
	TagNameJson    = "json"
	TagValueIgnore = "-"
)

type LabelObject interface{}

type MetricInfo struct {
	NameSpace  string
	SubSystem  string
	MetricName string
	Help       string
}

func NewMetricInfo(strNameSpace, strSubSystem, strMetricName, strHelp string) *MetricInfo {

	m := &MetricInfo{
		NameSpace:  strNameSpace,
		SubSystem:  strSubSystem,
		MetricName: strMetricName,
		Help:       strHelp,
	}
	return m
}

func (m *MetricInfo) Name() string {
	return fmt.Sprintf("%s_%s_%s", m.NameSpace, m.SubSystem, m.MetricName)
}

func (m *MetricInfo) newPrometheusDesc(labels ...string) *prometheus.Desc {
	strFQName := prometheus.BuildFQName(m.NameSpace, m.SubSystem, m.MetricName)
	return prometheus.NewDesc(strFQName, m.Help, labels, nil)
}

func (m *MetricInfo) newConstMetric(valueType prometheus.ValueType, obj LabelObject, value float64) Metrics {
	labelNames, labelValues := parseLabels(obj)
	desc := m.newPrometheusDesc(labelNames...)
	metric := prometheus.MustNewConstMetric(desc, valueType, value, labelValues...)
	return newConstMetric(m.Name(), metric, labelNames, labelValues)
}

func (m *MetricInfo) NewConstMetricGauge(obj LabelObject, value float64) Metrics {
	return m.newConstMetric(prometheus.GaugeValue, obj, value)
}

func (m *MetricInfo) NewConstMetricCounter(obj LabelObject, value float64) Metrics {
	return m.newConstMetric(prometheus.CounterValue, obj, value)
}

func (m *MetricInfo) NewConstMetricUntyped(obj LabelObject, value float64) Metrics {
	return m.newConstMetric(prometheus.UntypedValue, obj, value)
}
