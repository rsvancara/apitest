package handlers

import (
	"fmt"
	"net/http"

	"rainier/internal/util"

	"github.com/flosch/pongo2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus Metrics
var (
	eventsTotalProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "events_total",
		Help: "The total number of processed events",
	})

	eventsSearchLayerFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "events_searchlayer_failed_total",
		Help: "The total number of processed events failed",
	})

	eventsLocationHelperFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "events_locationhelper_failed_total",
		Help: "The total number of processed events failed",
	})

	eventsSearchLayerSucceed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "events_searchlayer_succeed_total",
		Help: "The total number of processed events succeed",
	})

	eventsLocationHelperSucceed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "events_locationhelper_succeed_total",
		Help: "The total number of processed events succeed",
	})

	eventSearchLayertLatency = promauto.NewSummary(prometheus.SummaryOpts{
		Name:       "latency_searchlayer_seconds",
		Help:       "Upstream latency in seconds.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})

	eventLocationHelperLatency = promauto.NewSummary(prometheus.SummaryOpts{
		Name:       "latency_locationhelper_seconds",
		Help:       "Upstream latency in seconds.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
)

func (ctx *HTTPHandlerContext) IndexHandler(w http.ResponseWriter, r *http.Request) {

	template, err := util.SiteTemplate("/main.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("error loading template with error: %s\n", err)
	}
	tmpl := pongo2.Must(pongo2.FromFile(template))

	out, err := tmpl.Execute(pongo2.Context{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("error loading template with error: %s\n", err)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, out)
}

// Reverse Phone Search
func (ctx *HTTPHandlerContext) TestHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"test\":\"success\"}")
}

func (ctx *HTTPHandlerContext) ErrorHandler(w http.ResponseWriter, r *http.Request) {

	template, err := util.SiteTemplate("/error.html")
	tmpl := pongo2.Must(pongo2.FromFile(template))

	out, err := tmpl.Execute(pongo2.Context{
		"title":    "Index",
		"greating": "Hello",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("error loading template with error: %s\n", err)
	}
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, out)
}

func (ctx *HTTPHandlerContext) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "working")
}

func NotFound(w http.ResponseWriter, r *http.Request) { // a * before http.Request
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\"status\":\"not found\"}")

}
