package prometheus

import (
	"strings"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/prom2json"
)

type Metric struct {
	Name      string            `json:"name"`
	Help      string            `json:"help,omitempty"`
	Type      string            `json:"type,omitempty"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
	Labels    map[string]string `json:"tags,omitempty"`
}

func flatMetricFamilies(mfs []*dto.MetricFamily) []*Metric {
	var metrics []*Metric
	for _, mf := range mfs {
		if ms := flatMetricFamily(mf); len(ms) != 0 {
			metrics = append(metrics, ms...)
		}
	}
	timestamp := GetTimestampMS()
	for _, m := range metrics {
		m.Timestamp = timestamp
	}
	return metrics
}

func flatMetricFamily(mf *dto.MetricFamily) []*Metric {
	var metrics []*Metric
	f := prom2json.NewFamily(mf)
	for _, m := range f.Metrics {
		switch v := m.(type) {
		case prom2json.Metric:
			metrics = append(metrics, flatMetric(&v, f))
		case prom2json.Summary:
			metrics = append(metrics, flatSummaryMetric(&v, f)...)
		case prom2json.Histogram:
			metrics = append(metrics, flatHistogramMetric(&v, f)...)
		}
	}
	return metrics
}

func flatMetric(m *prom2json.Metric, mf *prom2json.Family) *Metric {
	v, err := StringToFloat(m.Value)
	if err != nil {
		v = -1
	}
	return &Metric{
		Name:   mf.Name,
		Help:   mf.Help,
		Type:   mf.Type,
		Value:  v,
		Labels: m.Labels,
	}
}

func flatCountMetric(cmi interface{}, mf *prom2json.Family) *Metric {
	var value string
	var labels map[string]string

	switch v := cmi.(type) {
	case *prom2json.Summary:
		value = v.Count
		labels = v.Labels
	case *prom2json.Histogram:
		value = v.Count
		labels = v.Labels
	default:
		return nil
	}

	v, err := StringToFloat(value)
	if err != nil {
		v = -1
	}
	return &Metric{
		Name:   joinMetricName(mf.Name, "count"),
		Help:   mf.Help,
		Type:   mf.Type,
		Value:  v,
		Labels: labels,
	}
}

func flatSummaryMetric(summary *prom2json.Summary, mf *prom2json.Family) []*Metric {
	var metrics []*Metric
	if summary.Count != "" {
		metrics = append(metrics, flatCountMetric(summary, mf))
	}
	if summary.Sum != "" {
		metrics = append(metrics, flatSumMetric(summary, mf))
	}
	for qk, qv := range summary.Quantiles {
		v, err := StringToFloat(qv)
		if err != nil {
			v = -1
		}
		metric := &Metric{
			Name:  mf.Name,
			Help:  mf.Help,
			Type:  mf.Type,
			Value: v,
		}
		labels := make(map[string]string)
		for lk, lv := range summary.Labels {
			labels[lk] = lv
		}
		labels["quantile"] = qk
		metric.Labels = labels
		metrics = append(metrics, metric)
	}
	return metrics
}

func flatHistogramMetric(histogram *prom2json.Histogram, mf *prom2json.Family) []*Metric {
	var metrics []*Metric
	if histogram.Count != "" {
		metrics = append(metrics, flatCountMetric(histogram, mf))
	}
	if histogram.Sum != "" {
		metrics = append(metrics, flatSumMetric(histogram, mf))
	}
	for bk, bv := range histogram.Buckets {
		v, err := StringToFloat(bv)
		if err != nil {
			v = -1
		}
		metric := &Metric{
			Name:  joinMetricName(mf.Name, "bucket"),
			Help:  mf.Help,
			Type:  mf.Type,
			Value: v,
		}
		labels := make(map[string]string)
		for lk, lv := range histogram.Labels {
			labels[lk] = lv
		}
		labels["le"] = bk
		metric.Labels = labels
		metrics = append(metrics, metric)
	}
	return metrics
}

func flatSumMetric(smi interface{}, mf *prom2json.Family) *Metric {
	var value string
	var labels map[string]string

	switch v := smi.(type) {
	case *prom2json.Summary:
		value = v.Sum
		labels = v.Labels
	case *prom2json.Histogram:
		value = v.Sum
		labels = v.Labels
	default:
		return nil
	}

	v, err := StringToFloat(value)
	if err != nil {
		v = -1
	}
	return &Metric{
		Name:   joinMetricName(mf.Name, "sum"),
		Help:   mf.Help,
		Type:   mf.Type,
		Value:  v,
		Labels: labels,
	}
}

func joinMetricName(s ...string) string {
	return strings.Join(s, "_")
}
