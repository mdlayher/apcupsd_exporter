package apcupsdexporter

import (
	"log"

	"github.com/mdlayher/apcupsd"
	"github.com/prometheus/client_golang/prometheus"
)

var _ StatusSource = &apcupsd.Client{}

// A StatusSource is a type which can retrieve UPS status information from
// apcupsd.  It is implemented by *apcupsd.Client.
type StatusSource interface {
	Status() (*apcupsd.Status, error)
}

// A UPSCollector is a Prometheus collector for metrics regarding an APC UPS.
type UPSCollector struct {
	UPSLoadPercent                      *prometheus.Desc
	BatteryChargePercent                *prometheus.Desc
	LineVolts                           *prometheus.Desc
	LineNominalVolts                    *prometheus.Desc
	BatteryVolts                        *prometheus.Desc
	BatteryNominalVolts                 *prometheus.Desc
	BatteryNumberTransfersTotal         *prometheus.Desc
	BatteryTimeLeftSeconds              *prometheus.Desc
	BatteryTimeOnSeconds                *prometheus.Desc
	BatteryCumulativeTimeOnSecondsTotal *prometheus.Desc

	ss StatusSource
}

var _ prometheus.Collector = &UPSCollector{}

// NewUPSCollector creates a new UPSCollector.
func NewUPSCollector(ss StatusSource) *UPSCollector {
	var (
		labels = []string{"hostname", "ups_name", "model"}
	)

	return &UPSCollector{
		UPSLoadPercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "ups_load_percent"),
			"Current UPS load percentage.",
			labels,
			nil,
		),

		BatteryChargePercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_charge_percent"),
			"Current UPS battery charge percentage.",
			labels,
			nil,
		),

		LineVolts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "line_volts"),
			"Current AC input line voltage.",
			labels,
			nil,
		),

		LineNominalVolts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "line_nominal_volts"),
			"Nominal AC input line voltage.",
			labels,
			nil,
		),

		BatteryVolts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_volts"),
			"Current UPS battery voltage.",
			labels,
			nil,
		),

		BatteryNominalVolts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_nominal_volts"),
			"Nominal UPS battery voltage.",
			labels,
			nil,
		),

		BatteryNumberTransfersTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_number_transfers_total"),
			"Total number of transfers to UPS battery power.",
			labels,
			nil,
		),

		BatteryTimeLeftSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_time_left_seconds"),
			"Number of seconds remaining of UPS battery power.",
			labels,
			nil,
		),

		BatteryTimeOnSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_time_on_seconds"),
			"Number of seconds the UPS has been providing battery power due to an AC input line outage.",
			labels,
			nil,
		),

		BatteryCumulativeTimeOnSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_cumulative_time_on_seconds_total"),
			"Total number of seconds the UPS has provided battery power due to AC input line outages.",
			labels,
			nil,
		),

		ss: ss,
	}
}

// collect begins a metrics collection task for all metrics related to an APC
// UPS.
func (c *UPSCollector) collect(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	s, err := c.ss.Status()
	if err != nil {
		return c.BatteryVolts, err
	}

	labels := []string{
		s.Hostname,
		s.UPSName,
		s.Model,
	}

	ch <- prometheus.MustNewConstMetric(
		c.UPSLoadPercent,
		prometheus.GaugeValue,
		s.LoadPercent,
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryChargePercent,
		prometheus.GaugeValue,
		s.BatteryChargePercent,
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.LineVolts,
		prometheus.GaugeValue,
		s.LineVoltage,
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.LineNominalVolts,
		prometheus.GaugeValue,
		s.NominalInputVoltage,
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryVolts,
		prometheus.GaugeValue,
		s.BatteryVoltage,
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryNominalVolts,
		prometheus.GaugeValue,
		s.NominalBatteryVoltage,
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryNumberTransfersTotal,
		prometheus.CounterValue,
		float64(s.NumberTransfers),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryTimeLeftSeconds,
		prometheus.GaugeValue,
		s.TimeLeft.Seconds(),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryTimeOnSeconds,
		prometheus.GaugeValue,
		s.TimeOnBattery.Seconds(),
		labels...,
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryCumulativeTimeOnSecondsTotal,
		prometheus.CounterValue,
		s.CumulativeTimeOnBattery.Seconds(),
		labels...,
	)

	return nil, nil
}

// Describe sends the descriptors of each metric over to the provided channel.
// The corresponding metric values are sent separately.
func (c *UPSCollector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.UPSLoadPercent,
		c.BatteryChargePercent,
		c.LineVolts,
		c.LineNominalVolts,
		c.BatteryVolts,
		c.BatteryNominalVolts,
		c.BatteryNumberTransfersTotal,
		c.BatteryTimeLeftSeconds,
		c.BatteryTimeOnSeconds,
		c.BatteryCumulativeTimeOnSecondsTotal,
	}

	for _, d := range ds {
		ch <- d
	}
}

// Collect sends the metric values for each metric created by the UPSCollector
// to the provided prometheus Metric channel.
func (c *UPSCollector) Collect(ch chan<- prometheus.Metric) {
	if desc, err := c.collect(ch); err != nil {
		log.Printf("[ERROR] failed collecting UPS metric %v: %v", desc, err)
		ch <- prometheus.NewInvalidMetric(desc, err)
		return
	}
}
