package main

import (
	"flag"
	"fmt"
	"net/http"

	"crypto/sha256"
	"crypto/subtle"
	"static-metrics-exporter/internal/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var expectedUsername string
var expectedPassword string

var gaugeMetric prometheus.Gauge = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "example_gauge",
	Help: "This is example gauge metric",
})

func init() {
	prometheus.MustRegister(gaugeMetric)
	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
}

func main() {
	configPath := flag.String("config", "", "Path to YAML config file")
	flag.Parse()

	cfg := config.MustLoad(*configPath)
	for username, password := range cfg.BasicAuth {
		expectedUsername = username
		expectedPassword = password
	}
	fmt.Println("Started http server...")
	http.Handle("/metrics", basicAuthMiddleware(promhttp.Handler()))
	http.ListenAndServe(":2121", nil)
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
