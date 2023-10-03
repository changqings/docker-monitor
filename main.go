// package main

// import (
// 	"docker-client/status"
// 	"fmt"
// 	"log"
// 	"os"
// )

// func main() {

// 	m, err := status.GetMemUsage()
// 	if err != nil {
// 		log.Printf("GetMemUsage err: %v\n", err)
// 		os.Exit(1)
// 	}

// 	for _, v := range m {
// 		fmt.Printf("容器 name = %s ,容器 id = %s 的内存使用率为 %.2f%%\n", v.Id[:12], v.Name, v.Usage)
// 	}

// }

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
	go metrics.RecordMetrics(interval, wg)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":18085", nil)
}
