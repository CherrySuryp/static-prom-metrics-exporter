package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"crypto/sha256"
	"crypto/subtle"
	"static-metrics-exporter/internal/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var expectedUsername string
var expectedPassword string

func init() {
	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}

func main() {
	configPath := flag.String("config", "", "Path to YAML config file")
	flag.Parse()

	cfg := config.MustLoad(*configPath)
	for username, password := range cfg.Server.BasicAuth {
		expectedUsername = username
		expectedPassword = password
	}

	MustRegisterStaticMetrics(cfg.StaticMetrics)

	http.Handle("/metrics", basicAuthMiddleware(promhttp.Handler()))

	fmt.Printf("Starting http server on port: %s\n", cfg.Server.Port)

	var httpAddress string = fmt.Sprintf(":%s", cfg.Server.Port)
	if cfg.Server.TlsCrt != "" && cfg.Server.TlsKey != "" {
		if _, err := os.Stat(cfg.Server.TlsCrt); os.IsNotExist(err) {
			log.Fatal("Couldn't open TLS Certificate File")
		}
		if _, err := os.Stat(cfg.Server.TlsKey); os.IsNotExist(err) {
			log.Fatal("Couldn't open TLS Certificate File")
		}
		fmt.Println("TLS enabled")
		http.ListenAndServeTLS(httpAddress, cfg.Server.TlsCrt, cfg.Server.TlsKey, nil)
	}
	http.ListenAndServe(httpAddress, nil)
}

func MustRegisterStaticMetrics(metrics []config.StaticMetric) {
	fmt.Println("Initializing metrics...")
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

func basicAuthMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		username, password, ok := r.BasicAuth()
		if ok {
			presentUsernameHash := sha256.Sum256([]byte(username))
			presentPasswordHash := sha256.Sum256([]byte(password))
			
			expectedUsernameHash := sha256.Sum256([]byte(expectedUsername))
			expectedPasswordHash := sha256.Sum256([]byte(expectedPassword))

			usernameMatch := subtle.ConstantTimeCompare(presentUsernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(presentPasswordHash[:], expectedPasswordHash[:]) == 1
			
			if usernameMatch && passwordMatch{
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
