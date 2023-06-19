package main

import (
	"github.com/civet148/log"
	"github.com/civet148/promec"
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
	log.SetLevel("debug")
	client := promec.NewPromeClient()
	CommitMetrics(client)
	if err := client.Listen(":8088"); err != nil {
		log.Errorf("listen error: %s", err)
	}
}

func CommitMetrics(c *promec.PromeClient) {
	/*
		$ curl 127.0.0.1:8088/metrics | grep golang
		  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
		                                 Dload  Upload   Total   Spent    Left  Speed
		100  6151    0  6151    0     0  6006k      0 --:--:-- --:--:-- --:--:-- 6006k
		# HELP golang_prometheus_performance
		# TYPE golang_prometheus_performance gauge
		golang_prometheus_performance{cpu_usage="0.6",memory_usage="0.913",program="program1"} 1
		golang_prometheus_performance{cpu_usage="0.7",memory_usage="0.03",program="program2"} 1
	*/
	mi := promec.NewMetricInfo(NameSpace, SubSystem, MetricName, HelpName)
	gaugeLabels := make(map[*GaugeLabels]float64)
	gls := &GaugeLabels{
		Program:     "program1",
		CpuUsage:    0.6,
		MemoryUsage: 0.913,
	}

	gls2 := &GaugeLabels{
		Program:     "program2",
		CpuUsage:    0.7,
		MemoryUsage: 0.03,
	}
	gaugeLabels[gls] = 1
	gaugeLabels[gls2] = 1
	var metrics []promec.Metrics
	for label, value := range gaugeLabels {
		g := mi.NewConstMetricGauge(label, value)
		metrics = append(metrics, g)
	}
	c.WriteMetrics(metrics...)
}
