package metrics

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Registry collects counters and histograms and exposes them in
// Prometheus text exposition format at /metrics.
type Registry struct {
	mu         sync.Mutex
	counters   map[string]*counter
	histograms map[string]*histogram
	startedAt  time.Time
}

type counter struct {
	help   string
	labels map[string]uint64 // label-value -> count
}

type histogram struct {
	help    string
	buckets []float64
	counts  []uint64
	sum     float64
	total   uint64
}

func NewRegistry() *Registry {
	return &Registry{
		counters:   make(map[string]*counter),
		histograms: make(map[string]*histogram),
		startedAt:  time.Now(),
	}
}

// IncCounter increments a named counter by one.
func (r *Registry) IncCounter(name, help, labelValue string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	c, ok := r.counters[name]
	if !ok {
		c = &counter{help: help, labels: make(map[string]uint64)}
		r.counters[name] = c
	}
	c.labels[labelValue]++
}

var defaultBuckets = []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}

// Observe records a duration in seconds into a named histogram.
func (r *Registry) Observe(name, help string, seconds float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	h, ok := r.histograms[name]
	if !ok {
		h = &histogram{help: help, buckets: defaultBuckets, counts: make([]uint64, len(defaultBuckets))}
		r.histograms[name] = h
	}
	for i, upper := range h.buckets {
		if seconds <= upper {
			h.counts[i]++
		}
	}
	h.sum += seconds
	h.total++
}

// Render produces the Prometheus text exposition payload.
func (r *Registry) Render() string {
	r.mu.Lock()
	defer r.mu.Unlock()

	var b strings.Builder

	fmt.Fprintf(&b, "# HELP process_uptime_seconds Seconds since process start.\n")
	fmt.Fprintf(&b, "# TYPE process_uptime_seconds gauge\n")
	fmt.Fprintf(&b, "process_uptime_seconds %d\n", int64(time.Since(r.startedAt).Seconds()))

	names := make([]string, 0, len(r.counters))
	for name := range r.counters {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		c := r.counters[name]
		fmt.Fprintf(&b, "# HELP %s %s\n", name, c.help)
		fmt.Fprintf(&b, "# TYPE %s counter\n", name)
		labelValues := make([]string, 0, len(c.labels))
		for lv := range c.labels {
			labelValues = append(labelValues, lv)
		}
		sort.Strings(labelValues)
		for _, lv := range labelValues {
			fmt.Fprintf(&b, "%s{result=%q} %d\n", name, lv, c.labels[lv])
		}
	}

	hnames := make([]string, 0, len(r.histograms))
	for name := range r.histograms {
		hnames = append(hnames, name)
	}
	sort.Strings(hnames)
	for _, name := range hnames {
		h := r.histograms[name]
		fmt.Fprintf(&b, "# HELP %s %s\n", name, h.help)
		fmt.Fprintf(&b, "# TYPE %s histogram\n", name)
		for i, upper := range h.buckets {
			fmt.Fprintf(&b, "%s_bucket{le=%q} %d\n", name, fmt.Sprintf("%g", upper), h.counts[i])
		}
		fmt.Fprintf(&b, "%s_bucket{le=\"+Inf\"} %d\n", name, h.total)
		fmt.Fprintf(&b, "%s_sum %g\n", name, h.sum)
		fmt.Fprintf(&b, "%s_count %d\n", name, h.total)
	}

	return b.String()
}

// Middleware instruments HTTP requests with request count and latency,
// then registers GET /metrics on the router.
func Middleware(reg *Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()

		status := fmt.Sprintf("%d", c.Writer.Status())
		reg.IncCounter("http_requests_total", "Total HTTP requests.", status)
		reg.Observe("http_request_duration_seconds", "HTTP request latency in seconds.", time.Since(start).Seconds())
	}
}

// Handler serves the /metrics endpoint.
func Handler(reg *Registry) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain; version=0.0.4; charset=utf-8", []byte(reg.Render()))
	}
}
