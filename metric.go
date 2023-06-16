package promec

import (
	"fmt"
	"github.com/civet148/log"
	"github.com/prometheus/client_golang/prometheus"
	"reflect"
)

const (
	TagNameLabel   = "label"
	TagNameJson    = "json"
	TagValueIgnore = "-"
)

type MetricStruct interface{}

type Metric struct {
	NameSpace  string
	SubSystem  string
	MetricName string
	Help       string
	FQName     string
}

func NewMetric(strNameSpace, strSubSystem, strMetricName, strHelp string) *Metric {
	strFQName := prometheus.BuildFQName(strNameSpace, strSubSystem, strMetricName)
	m := &Metric{
		NameSpace:  strNameSpace,
		SubSystem:  strSubSystem,
		MetricName: strMetricName,
		Help:       strHelp,
		FQName:     strFQName,
	}
	return m
}

func (m *Metric) newPrometheusDesc(labels ...string) *prometheus.Desc {
	return prometheus.NewDesc(m.FQName, m.Help, labels, nil)
}

func (m *Metric) NewConstMetricGauge(obj MetricStruct, value float64) prometheus.Metric {
	labels, values := m.parseLabels(obj)
	if len(labels) == 0 {
		return nil
	}
	desc := m.newPrometheusDesc(labels...)
	return prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value, values...)
}

//
//func (m *Metric) NewConstMetricCounter(obj MetricStruct, value float64) prometheus.Metric {
//	labels, values := m.parseLabels(obj)
//	if len(labels) == 0 {
//		return nil
//	}
//	desc := m.newPrometheusDesc(labels...)
//	return prometheus.MustNewConstMetric(desc, prometheus.CounterValue, value, values...)
//}
//
//func (m *Metric) NewConstMetricUntyped(obj MetricStruct, value float64) prometheus.Metric {
//	labels, values := m.parseLabels(obj)
//	if len(labels) == 0 {
//		return nil
//	}
//	desc := m.newPrometheusDesc(labels...)
//	return prometheus.MustNewConstMetric(desc, prometheus.UntypedValue, value, values...)
//}

func (m *Metric) parseLabels(obj MetricStruct) (labels, values []string) {
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	for {
		if typ.Kind() != reflect.Ptr { // pointer type
			break
		}
		typ = typ.Elem()
		val = val.Elem()
	}

	kind := typ.Kind()
	switch kind {
	case reflect.Struct:
		{
			return m.parseStructFields(typ, val)
		}
	//case reflect.Slice:
	//	{
	//		typ = val.Type().Elem()
	//		val = reflect.New(typ).Elem()
	//		m.parseStructFields(typ, val)
	//	}
	default:
		log.Panic("object kind [%v] not support yet", typ.Kind())
	}
	return
}

// parse struct fields
func (m *Metric) parseStructFields(typ reflect.Type, val reflect.Value) (labels, values []string) {
	kind := typ.Kind()
	if kind == reflect.Struct {
		NumField := val.NumField()
		for i := 0; i < NumField; i++ {
			typField := typ.Field(i)
			valField := val.Field(i)

			if typField.Type.Kind() == reflect.Ptr {
				typField.Type = typField.Type.Elem()
				valField = valField.Elem()
			}
			if !valField.IsValid() || !valField.CanInterface() {
				continue
			}
			strTagVal, ignore := m.getTagValue(typField)
			if ignore {
				continue
			}
			if strTagVal != "" {
				labels = append(labels, strTagVal)
				values = append(values, fmt.Sprintf("%v", valField.Interface()))
			}
		}
	}
	return
}

// get struct field's tag value
func (m *Metric) getTagValue(sf reflect.StructField) (strValue string, ignore bool) {
	strValue = sf.Tag.Get(TagNameLabel)
	if strValue == TagValueIgnore {
		return "", true
	}
	if strValue == "" {
		strValue = sf.Tag.Get(TagNameJson)
		if strValue == TagValueIgnore {
			return "", true
		}
	}
	return
}
