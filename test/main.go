package main

import (
	"github.com/civet148/log"
	"github.com/civet148/promec"
	"math/rand"
	"time"
)

const (
	NameSpace  = "golang"
	SubSystem  = "prometheus"
	MetricName = "performance"
	HelpName   = ""
)

type GaugeLabels struct {
	Program     string  `json:"program"`
	Timestamp   int64   `json:"timestamp"`
	CpuUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}

func main() {
	client := promec.NewPromeClient()
	go MetricsTimerTask(client)
	if err := client.Listen(":8088"); err != nil {
		log.Errorf("listen error: %s", err)
	}
}

func MetricsTimerTask(c *promec.PromeClient) {
	ticker := time.NewTicker(5 * time.Second)
	metric := promec.NewMetric(NameSpace, SubSystem, MetricName, HelpName)

	for {
		select {
		case <-ticker.C:
			gls := &GaugeLabels{
				Program:     "promec",
				Timestamp:   time.Now().Unix(),
				CpuUsage:    0.8,
				MemoryUsage: 0.913,
			}
			value := rand.Int31n(100) //random integer
			v := metric.NewConstMetricGauge(gls, float64(value))
			c.WriteMetric(v)
			/*
					$ curl 127.0.0.1:8088/metrics| grep -vE "TYPE|HELP" | grep golang
					  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
					                                 Dload  Upload   Total   Spent    Left  Speed
					 100  6100    0  6100    0     0  1985k      0 --:--:-- --:--:-- --:--:-- 1985k
				     golang_prometheus_performance{cpu_usage="0.8",memory_usage="0.913",program="promec",timestamp="1686901954"} 47
			*/
		}
	}
}
