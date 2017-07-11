package go_metrics

import (
	"fmt"
	"github.com/moira-alert/moira-alert"
	"github.com/moira-alert/moira-alert/metrics/graphite"
	"github.com/rcrowley/go-metrics"
	"net"
	"os"
	"strings"
	"time"

	goMetricsGraphite "github.com/cyberdelia/go-metrics-graphite"
)

func Init(metric *graphite.NotifierMetrics, logger moira_alert.Logger) {
	config := (*metric).Config

	uri := config.URI
	prefix := config.Prefix
	interval := config.Interval

	if uri != "" {
		address, err := net.ResolveTCPAddr("tcp", uri)
		if err != nil {
			logger.Errorf("Can not resolve graphiteURI %s: %s", uri, err)
			return
		}
		hostname, err := os.Hostname()
		if err != nil {
			logger.Errorf("Can not get OS hostname: %s", err)
			return
		}
		shortName := strings.Split(hostname, ".")[0]
		go goMetricsGraphite.Graphite(metrics.DefaultRegistry, time.Duration(interval)*time.Second, fmt.Sprintf("%s.notifier.%s", prefix, shortName), address)
	}
}

func NewRegisteredMeter(name string) metrics.Meter {
	return metrics.NewRegisteredMeter("events.received", metrics.DefaultRegistry)
}

func ConfigureNotifierMetrics(config graphite.Config) graphite.NotifierMetrics {
	graphite.NotifierMetric = graphite.NotifierMetrics{
		Config:                 config,
		EventsReceived:         NewRegisteredMeter("events.received"),
		EventsMalformed:        NewRegisteredMeter("events.malformed"),
		EventsProcessingFailed: NewRegisteredMeter("events.failed"),
		SubsMalformed:          NewRegisteredMeter("subs.malformed"),
		SendingFailed:          NewRegisteredMeter("sending.failed"),
		SendersOkMetrics:       &MetricsMap{make(map[string]metrics.Meter)},
		SendersFailedMetrics:   &MetricsMap{make(map[string]metrics.Meter)},
	}
	return graphite.NotifierMetric
}