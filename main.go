package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	constValueOne = promauto.NewCounter(prometheus.CounterOpts{
		Name: "promtest_const_value_one_total",
		Help: "The constant number 1.0",
	})
)

func healthz(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handling healthz: %s", r.URL.Path[1:])
	fmt.Fprintf(w, "%s", r.URL.Path[1:])
}

func metricsWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Before Handling metrics request")
		h.ServeHTTP(w, r) // call original
		log.Println("After Handling metrics request")
	})
}

func main() {
	constValueOne.Inc()

	http.Handle("/metrics", metricsWrapper(promhttp.Handler()))
	http.HandleFunc("/healthz/readieness", healthz)
	http.HandleFunc("/healthz/liveness", healthz)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
