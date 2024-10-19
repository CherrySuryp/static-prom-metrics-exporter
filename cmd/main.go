package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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
	configPath := flag.String("config", "config.yml", "Path to YAML config file")
	port := flag.Int("port", 9002, "Port to listen on")
	flag.Parse()

	cfg := config.MustLoad(*configPath)
	for username, password := range cfg.Server.BasicAuth {
		expectedUsername = username
		expectedPassword = password
	}

	MustRegisterStaticMetrics(cfg.StaticMetrics)

	http.Handle("/metrics", BasicAuthMiddleware(promhttp.Handler()))

	startServer(fmt.Sprintf(":%d", *port), cfg.Server.TlsCrt, cfg.Server.TlsKey)
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

func MustRegisterStaticMetrics(metrics []config.StaticMetric) {
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

func BasicAuthMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			presentUsernameHash := sha256.Sum256([]byte(username))
			presentPasswordHash := sha256.Sum256([]byte(password))

			expectedUsernameHash := sha256.Sum256([]byte(expectedUsername))
			expectedPasswordHash := sha256.Sum256([]byte(expectedPassword))

			usernameMatch := subtle.ConstantTimeCompare(presentUsernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(presentPasswordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
