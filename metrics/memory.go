package metrics

import (
	"docker-client/status"
	"fmt"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RunGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "container_memory_usage",
		Help: "container memory usage as percent",
	},
		[]string{
			"container_name",
			"container_id",
			"container_limit",
		},
	)
)

var (
	RunSummary = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "container_memory_usage_sumary",
		Help:       "container memory usage sumary with 0.5,0.9,0.99",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	},
		[]string{
			"container_name",
			"container_id",
			"container_limit",
		},
	)
)

func RecordMetrics(m status.MemoryStat, interval *int, wg *sync.WaitGroup) {
	RunGauge.WithLabelValues(m.Name, m.Id[:12], fmt.Sprintf("%dM", m.MemLimit/1024/1024)).Set(decimal(m.Usage))
	RunSummary.WithLabelValues(m.Name, m.Id[:12], fmt.Sprintf("%dM", m.MemLimit/1024/1024)).Observe(decimal(m.Usage))
	wg.Done()
}

func decimal(f float64) float64 {
	f, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", f), 64)
	return f
}
