package main

import (
	"github.com/civet148/log"
	"github.com/civet148/promec"
)

const (
	NameSpace           = "golang"
	SubSystem           = "prometheus"
	MetricNameGauge     = "gauge"
	MetricNameHistogram = "histogram"
	HelpName            = ""
)

type GaugeLabels struct {
	Program     string  `json:"program"`
	CpuUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}

type HistogramLabels struct {
	Code   int32  `json:"code"`
	Method string `json:"method"`
}

func main() {
	log.SetLevel("debug")
	client := promec.NewPromeClient()
	CommitConstMetrics(client)
	CommitConstHistogram(client)
	if err := client.Listen(":8088"); err != nil {
		log.Errorf("listen error: %s", err)
	}
}

func CommitConstMetrics(c *promec.PromeClient) {
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
	mi := promec.NewMetricInfo(NameSpace, SubSystem, MetricNameGauge, HelpName)
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
	for obj, value := range gaugeLabels {
		g := mi.NewConstMetricGauge(obj, value)
		metrics = append(metrics, g)
	}
	c.WriteMetrics(metrics...)
}

func CommitConstHistogram(c *promec.PromeClient) {
	/*
		$ curl 127.0.0.1:8088/metrics | grep golang | grep -vE "HELP"
		  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
		                                 Dload  Upload   Total   Spent    Left  Speed
		100  6699    0  6699    0     0  6541k      0 --:--:-- --:--:-- --:--:-- 6541k
		# TYPE golang_prometheus_histogram histogram
		golang_prometheus_histogram_bucket{code="200",method="POST",le="0.5"} 10
		golang_prometheus_histogram_bucket{code="200",method="POST",le="1"} 20
		golang_prometheus_histogram_bucket{code="200",method="POST",le="2.5"} 30
		golang_prometheus_histogram_bucket{code="200",method="POST",le="3"} 40
		golang_prometheus_histogram_bucket{code="200",method="POST",le="+Inf"} 100
		golang_prometheus_histogram_sum{code="200",method="POST"} 130.52
		golang_prometheus_histogram_count{code="200",method="POST"} 100
	*/
	mi := promec.NewMetricInfo(NameSpace, SubSystem, MetricNameHistogram, HelpName)
	var buckets = map[float64]uint64{
		0.5: 10,
		1.0: 20,
		2.5: 30,
		3.0: 40,
	}
	labelPOST := &HistogramLabels{
		Code:   200,
		Method: "POST",
	}
	histo := mi.NewConstHistogram(labelPOST, 100, 130.52, buckets)
	c.WriteMetrics(histo)

	//go func() {
	//	time.Sleep(15 * time.Second)
	//	histo.Update(200, 309.5, buckets)
	//}()
}
