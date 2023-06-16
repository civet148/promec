package promec

import (
	"github.com/civet148/log"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

const (
	DefaultMetricsPath = "/metrics"
)

type PromeClient struct {
	metrics chan prometheus.Metric
}

func NewPromeClient() *PromeClient {
	m := &PromeClient{
		metrics: make(chan prometheus.Metric, 1000),
	}
	err := prometheus.Register(m)
	if err != nil {
		log.Errorf(err.Error())
	}
	return m
}

// Describe DO NOT USE THIS FUNCTION
func (m *PromeClient) Describe(ch chan<- *prometheus.Desc) {
	mc := make(chan prometheus.Metric) //metric channel
	dc := make(chan struct{})          //done channel
	go func() {
		for metric := range mc {
			ch <- metric.Desc()
		}
		close(dc)
	}()
	m.Collect(mc)
	close(mc)
	<-dc
}

// Collect DO NOT USE THIS FUNCTION
func (m *PromeClient) Collect(ch chan<- prometheus.Metric) {
	select {
	case metric := <-m.metrics:
		ch <- metric
	default:

	}
}

// WriteMetric write metric to prometheus buffer
func (m *PromeClient) WriteMetric(metric prometheus.Metric) {
	m.metrics <- metric
}

// InitRouter init prometheus router
func (m *PromeClient) InitRouter(r gin.IRouter, relativePath string) {
	r.GET(relativePath, gin.WrapH(promhttp.Handler()))
}

// Listen start prometheus metrics server
// strListenAddr: net address to listen (eg. 0.0.0.0:8080)
// relativePath: relative path to prometheus metrics server (default /metrics)
func (m *PromeClient) Listen(strListenAddr string, relativePath ...string) (err error) {

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	strMetricsPath := DefaultMetricsPath
	if len(relativePath) > 0 {
		strMetricsPath = relativePath[0]
	}
	m.InitRouter(r, strMetricsPath)
	log.Infof("prometheus metrics server listening on %s/%s", strListenAddr, strMetricsPath)
	if err = http.ListenAndServe(strListenAddr, r); err != nil { //if everything is fine, it will block this routine
		log.Errorf("listen http address [%s] error [%s]\n", strListenAddr, err)
	}
	return nil
}
