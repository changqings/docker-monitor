package main

import (
	"docker-client/metrics"
	"flag"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	// go cmd.Run ping()
	interval := flag.Int("i", 10, "-i timeSecond, unit = s")
	flag.Parse()

	// remove GC go_xx
	prometheus.Unregister(collectors.NewGoCollector())
	// remove process_xx
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	// This is should be goroutine in loop
	go metrics.RecordMetrics(*interval)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":18085", nil)
}
