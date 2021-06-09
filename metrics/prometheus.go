package metrics

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	service    string
	serviceKey = "_service"
)

func init() {
	service = os.Getenv("SERVICE_NAME")
	if service == "" {
		service = "unknown"
	}
}

// Counter is the counter-type metrics.
type Counter struct {
	counterVec *prometheus.CounterVec
}

// NewCounter returns a instance of Counter.
func NewCounter(name string, labels []string) *Counter {
	c := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        name,
			Help:        "-",
			ConstLabels: map[string]string{serviceKey: service},
		},
		labels)
	prometheus.MustRegister(c)
	return &Counter{
		counterVec: c,
	}
}

// Add adds the given value to the Counter.
func (c *Counter) Add(value float64, labels map[string]string) error {
	counter, err := c.counterVec.GetMetricWith(prometheus.Labels(labels))
	if err != nil {
		return err
	}
	counter.Add(value)
	return nil
}

// Gauge is the gauge-type metrics.
type Gauge struct {
	gaugeVec *prometheus.GaugeVec
}

// NewGauge returns a instance of Gauge.
func NewGauge(name string, labels []string) *Gauge {
	g := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        name,
			Help:        "-",
			ConstLabels: map[string]string{serviceKey: service},
		},
		labels)
	prometheus.MustRegister(g)
	return &Gauge{
		gaugeVec: g,
	}
}

// Set sets the given value to the Gauge.
func (g *Gauge) Set(value float64, labels map[string]string) error {
	gauge, err := g.gaugeVec.GetMetricWith(prometheus.Labels(labels))
	if err != nil {
		return err
	}
	gauge.Set(value)
	return nil
}

// Summary is summary-type metrics.
type Summary struct {
	summaryVec *prometheus.SummaryVec
}

// NewSummary returns a instance of Summary.
func NewSummary(name string, labels []string) *Summary {
	s := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:        name,
			Help:        "-",
			ConstLabels: map[string]string{serviceKey: service},
		},
		labels)
	prometheus.MustRegister(s)
	return &Summary{
		summaryVec: s,
	}
}

// Observe adds the observations to the Summary.
func (s *Summary) Observe(value float64, labels map[string]string) error {
	summary, err := s.summaryVec.GetMetricWith(prometheus.Labels(labels))
	if err != nil {
		return err
	}
	summary.Observe(value)
	return nil
}

// Histogram is the histogram-type metrics.
type Histogram struct {
	histogramVec *prometheus.HistogramVec
}

// NewHistogram returns a instance of Histogram.
func NewHistogram(name string, labels []string) *Histogram {
	h := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        name,
			Help:        "-",
			ConstLabels: map[string]string{serviceKey: service},
		},
		labels)
	prometheus.MustRegister(h)
	return &Histogram{
		histogramVec: h,
	}
}

// Observe adds the observations to the Histogram.
func (h *Histogram) Observe(value float64, labels map[string]string) error {
	histogram, err := h.histogramVec.GetMetricWith(prometheus.Labels(labels))
	if err != nil {
		return err
	}
	histogram.Observe(value)
	return nil
}
