package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"static-metrics-exporter/internal/middleware"

	"static-metrics-exporter/internal/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var flags Flags

var basicAuth middleware.BasicAuth

type Flags struct {
	ConfigPath string
	Port       int
	TlsCrt     string
	TlsKey     string
}

func (f *Flags) Init() {
	flag.StringVar(&f.ConfigPath, "config", "config.yml", "Path to YAML config file")
	flag.IntVar(&f.Port, "port", 9002, "Port to listen on")
	flag.StringVar(&f.TlsCrt, "tls-crt", "", "Path to TLS Certificate File")
	flag.StringVar(&f.TlsKey, "tls-key", "", "Path to TLS Key File")
	flag.Parse()
}

func mustRegisterStaticMetrics(metrics []config.StaticMetric) {
	for _, metric := range metrics {
		gaugeMetric := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: metric.Name,
			Help: metric.Help,
		})
		prometheus.MustRegister(gaugeMetric)
		gaugeMetric.Set(float64(metric.Value))
		fmt.Printf("Registered metric \"%s\" with value: %d\n", metric.Name, metric.Value)
	}
}

func initHandlers() {
	http.Handle("/metrics", basicAuth.Middleware(promhttp.Handler()))
}

func startServer(httpAddress string, tlsCrt string, tlsKey string) {
	fmt.Printf("Starting listening on port: %s\n", httpAddress)

	if tlsCrt != "" && tlsKey != "" {
		if _, err := os.Stat(tlsCrt); os.IsNotExist(err) {
			log.Fatal("Couldn't open TLS Certificate File")
		}
		if _, err := os.Stat(tlsKey); os.IsNotExist(err) {
			log.Fatal("Couldn't open TLS Key File")
		}
		fmt.Println("TLS enabled")
		err := http.ListenAndServeTLS(httpAddress, tlsCrt, tlsKey, nil)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	fmt.Println("TLS disabled")
	err := http.ListenAndServe(httpAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func init() {
	flags.Init()
	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	prometheus.Unregister(collectors.NewBuildInfoCollector())
	initHandlers()
}

func main() {
	cfg := config.MustLoad(flags.ConfigPath)
	mustRegisterStaticMetrics(cfg.StaticMetrics)

	var username string
	var password string
	for user, pass := range cfg.Server.BasicAuth {
		username = user
		password = pass
	}
	basicAuth.SetCredentials(username, password)

	startServer(fmt.Sprintf(":%d", flags.Port), flags.TlsCrt, flags.TlsKey)
}
