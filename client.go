package promec

import (
	"github.com/civet148/log"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
)

const (
	DefaultMetricsPath = "/metrics"
)

type PromeClient struct {
	locker  sync.RWMutex
	metrics map[string]Metrics
}

func NewPromeClient() *PromeClient {
	m := &PromeClient{
		metrics: make(map[string]Metrics),
	}
	prometheus.MustRegister(m)
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
	m.locker.RLock()
	defer m.locker.RUnlock()
	for _, v := range m.metrics {
		ch <- v.Metric()
	}
}

// WriteMetric write metric to prometheus buffer
func (m *PromeClient) WriteMetrics(metrics ...Metrics) {
	m.locker.Lock()
	defer m.locker.Unlock()
	for _, v := range metrics {
		m.metrics[v.Key()] = v
	}
}

// CleanMetric clean metrics buffer
func (m *PromeClient) CleanMetrics() {
	m.locker.Lock()
	defer m.locker.Unlock()
	for _, v := range m.metrics {
		delete(m.metrics, v.Key())
	}
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
