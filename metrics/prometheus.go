package metrics

import (
	"net/http"
	"sync"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Prometheus wrapper containing utility methods for registering and retrieving different metrics.
type Prometheus struct {
	mu         *sync.RWMutex
	registerer prometheus.Registerer
	gauges     map[string]prometheus.Gauge
	counters   map[string]prometheus.Counter
	histograms map[string]prometheus.Histogram
	summaries  map[string]prometheus.Summary
}

// NewPrometheus returns a new instance of Prometheus.
func NewPrometheus(registerer prometheus.Registerer) *Prometheus {
	if registerer == nil {
		registerer = prometheus.DefaultRegisterer
	}
	return &Prometheus{
		mu:         new(sync.RWMutex),
		registerer: registerer,
		gauges:     make(map[string]prometheus.Gauge),
		counters:   make(map[string]prometheus.Counter),
		histograms: make(map[string]prometheus.Histogram),
		summaries:  make(map[string]prometheus.Summary),
	}
}

// Handler returns the Prometheus http handler.
func (p *Prometheus) Handler() http.Handler {
	return promhttp.Handler()
}

// RegisterGauges registers the provided gauge metrics to the Prometheus registerer.
func (p *Prometheus) RegisterGauges(opts ...prometheus.GaugeOpts) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, options := range opts {
		p.registerGaugeIfNotExists(options)
	}
}

// UnregisterGauges unregisters the provided gauge metrics from the Prometheus registerer.
func (p *Prometheus) UnregisterGauges(names ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, name := range names {
		p.unregisterGaugeIfExists(name)
	}
}

// GetGauge retrieves gauge metric by name
func (p *Prometheus) GetGauge(name string) (gauge prometheus.Gauge, exist bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if gauge, exist = p.gauges[name]; !exist {
		return nil, exist
	}

	return gauge, exist
}

// RegisterCounters registers the provided counter metrics to the Prometheus registerer.
func (p *Prometheus) RegisterCounters(opts ...prometheus.CounterOpts) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, options := range opts {
		p.registerCounterIfNotExists(options)
	}
}

// GetCounter retrieves counter metric by name
func (p *Prometheus) GetCounter(name string) (counter prometheus.Counter, exist bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if counter, exist = p.counters[name]; !exist {
		return nil, exist
	}

	return counter, exist
}

// RegisterHistograms registers the provided histogram metrics to the Prometheus registerer.
func (p *Prometheus) RegisterHistograms(opts ...prometheus.HistogramOpts) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, options := range opts {
		p.registerHistogramIfNotExists(options)
	}
}

// GetHistogram retrieves histogram metric by name
func (p *Prometheus) GetHistogram(name string) (histogram prometheus.Histogram, exist bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if histogram, exist = p.histograms[name]; !exist {
		return nil, exist
	}

	return histogram, exist
}

// RegisterSummaries registers the provided summary metrics to the Prometheus registerer.
func (p *Prometheus) RegisterSummaries(opts ...prometheus.SummaryOpts) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, options := range opts {
		p.registerSummaryIfNotExists(options)
	}
}

// GetSummary retrieves summary metric by name
func (p *Prometheus) GetSummary(name string) (summary prometheus.Summary, exist bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if summary, exist = p.summaries[name]; !exist {
		return nil, exist
	}

	return summary, exist
}

// registerGaugeIfNotExists registers single gauge metric if not exists
func (p *Prometheus) registerGaugeIfNotExists(opts prometheus.GaugeOpts) {
	if _, exist := p.gauges[opts.Name]; exist {
		log.Warnf("Gauge metric '%v' already exists.", opts.Name)
		return
	}

	log.Infof("Creating Gauge Metric '%v' ...", opts.Name)
	gauge := prometheus.NewGauge(opts)
	log.Infof("Gauge Metric '%v' successfully created! Labels: %p", opts.Name, opts.ConstLabels)

	log.Infof("Registering Gauge Metric '%v' ...", opts.Name)
	p.registerer.MustRegister(gauge)
	log.Infof("Gauge Metric '%v' successfully registered!", opts.Name)

	p.gauges[opts.Name] = gauge
}

// unregisterGaugeIfExists unregisters single gauge metric if exists
func (p *Prometheus) unregisterGaugeIfExists(name string) {
	var (
		gauge prometheus.Gauge
		ok    bool
	)

	if gauge, ok = p.gauges[name]; !ok {
		log.Warnf("Trying to delete non-existing Gauge gauge '%v'.", name)
		return
	}

	log.Infof("Unregistering Gauge Metric '%v' ...", name)
	ok = p.registerer.Unregister(gauge)
	if !ok {
		log.Errorf("Failed to unregister Gauge Metric '%v'.", name)
		return
	}
	delete(p.gauges, name)
	log.Infof("Gauge Metric '%v' successfully unregistered!", name)
}

// registerCounterIfNotExists registers single counter metric if not exists
func (p *Prometheus) registerCounterIfNotExists(opts prometheus.CounterOpts) {
	if _, exist := p.counters[opts.Name]; exist {
		log.Warnf("Counter metric '%v' already exists.", opts.Name)
		return
	}

	log.Infof("Creating Counter Metric '%v' ...", opts.Name)
	counter := prometheus.NewCounter(opts)
	log.Infof("Counter Metric '%v' successfully created! Labels: %p", opts.Name, opts.ConstLabels)

	log.Infof("Registering Counter Metric '%v' ...", opts.Name)
	p.registerer.MustRegister(counter)
	log.Infof("Counter Metric '%v' successfully registered!", opts.Name)

	p.counters[opts.Name] = counter
}

// unregisterCounterIfExists unregisters single counter metric if exists
func (p *Prometheus) unregisterCounterIfExists(name string) {
	var (
		counter prometheus.Counter
		ok      bool
	)

	if counter, ok = p.counters[name]; !ok {
		log.Warnf("Trying to delete non-existing Counter counter '%v'.", name)
		return
	}

	log.Infof("Unregistering Counter Metric '%v' ...", name)
	ok = p.registerer.Unregister(counter)
	if !ok {
		log.Errorf("Failed to unregister Counter Metric '%v'.", name)
		return
	}
	delete(p.counters, name)
	log.Infof("Counter Metric '%v' successfully unregistered!", name)
}

// registerHistogramIfNotExists registers single histogram metric if not exists
func (p *Prometheus) registerHistogramIfNotExists(opts prometheus.HistogramOpts) {
	if _, exist := p.histograms[opts.Name]; exist {
		log.Warnf("Histogram metric '%v' already exists.", opts.Name)
		return
	}

	log.Infof("Creating Histogram Metric '%v' ...", opts.Name)
	histogram := prometheus.NewHistogram(opts)
	log.Infof("Histogram Metric '%v' successfully created! Labels: %p", opts.Name, opts.ConstLabels)

	log.Infof("Registering Histogram Metric '%v' ...", opts.Name)
	p.registerer.MustRegister(histogram)
	log.Infof("Histogram Metric '%v' successfully registered!", opts.Name)

	p.histograms[opts.Name] = histogram
}

// unregisterHistogramIfExists unregisters single histogram metric if exists
func (p *Prometheus) unregisterHistogramIfExists(name string) {
	var (
		histogram prometheus.Histogram
		ok        bool
	)

	if histogram, ok = p.histograms[name]; !ok {
		log.Warnf("Trying to delete non-existing Histogram histogram '%v'.", name)
		return
	}

	log.Infof("Unregistering Histogram Metric '%v' ...", name)
	ok = p.registerer.Unregister(histogram)
	if !ok {
		log.Errorf("Failed to unregister Histogram Metric '%v'.", name)
		return
	}
	delete(p.histograms, name)
	log.Infof("Histogram Metric '%v' successfully unregistered!", name)
}

// registerSummaryIfNotExists registers single summary metric if not exists
func (p *Prometheus) registerSummaryIfNotExists(opts prometheus.SummaryOpts) {
	if _, exist := p.summaries[opts.Name]; exist {
		log.Warnf("Summary metric '%v' already exists.", opts.Name)
		return
	}

	log.Infof("Creating Summary Metric '%v' ...", opts.Name)
	summary := prometheus.NewSummary(opts)
	log.Infof("Summary Metric '%v' successfully created! Labels: %p", opts.Name, opts.ConstLabels)

	log.Infof("Registering Summary Metric '%v' ...", opts.Name)
	p.registerer.MustRegister(summary)
	log.Infof("Summary Metric '%v' successfully registered!", opts.Name)

	p.summaries[opts.Name] = summary
}

// unregisterSummaryIfExists unregisters single summary metric if exists
func (p *Prometheus) unregisterSummaryIfExists(name string) {
	var (
		summary prometheus.Summary
		ok      bool
	)

	if summary, ok = p.summaries[name]; !ok {
		log.Warnf("Trying to delete non-existing Summary summary '%v'.", name)
		return
	}

	log.Infof("Unregistering Summary Metric '%v' ...", name)
	ok = p.registerer.Unregister(summary)
	if !ok {
		log.Errorf("Failed to unregister Summary Metric '%v'.", name)
		return
	}
	delete(p.summaries, name)
	log.Infof("Summary Metric '%v' successfully unregistered!", name)
}
