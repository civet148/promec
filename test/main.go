package main

import (
	"github.com/civet148/log"
	"github.com/civet148/promec"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	NameSpace  = "golang"
	SubSystem  = "prometheus"
	MetricName = "performance"
	HelpName   = ""
)

type GaugeLabels struct {
	Program     string  `json:"program"`
	CpuUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}

func main() {
	client := promec.NewPromeClient()
	CommitMetrics(client)
	if err := client.Listen(":8088"); err != nil {
		log.Errorf("listen error: %s", err)
	}
}

func CommitMetrics(c *promec.PromeClient) {
	mi := promec.NewMetricInfo(NameSpace, SubSystem, MetricName, HelpName)
	gaugeLabels := make(map[*GaugeLabels]float64)
	gls := &GaugeLabels{
		Program:     "program1",
		CpuUsage:    0.6,
		MemoryUsage: 0.913,
	}

	gls2 := &GaugeLabels{
		Program:     "program2",
		CpuUsage:    0.48,
		MemoryUsage: 0.609,
	}
	gaugeLabels[gls] = 1
	gaugeLabels[gls2] = 1
	var metrics []prometheus.Metric
	for label, value := range gaugeLabels {
		g := mi.NewConstMetricGauge(label, value)
		metrics = append(metrics, g)
	}
	c.WriteMetric(mi.Key(), metrics...)
}
