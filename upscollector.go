package apcupsdexporter

import (
	"log"
	"strings"
	"time"

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
	LastTransferOnBattery               *prometheus.Desc
	LastTransferOffBattery              *prometheus.Desc
	LastSelftest                        *prometheus.Desc
	NominalPowerWatts                   *prometheus.Desc
	UPSStatus                           *prometheus.Desc
	UPSInfo                             *prometheus.Desc

	ss StatusSource
}

var _ prometheus.Collector = &UPSCollector{}

// NewUPSCollector creates a new UPSCollector.
func NewUPSCollector(ss StatusSource) *UPSCollector {
	var (
		infoLabels = []string{"hostname", "ups_name", "model"}
	)

	return &UPSCollector{
		UPSLoadPercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "ups_load_percent"),
			"Current UPS load percentage.",
			nil,
			nil,
		),

		BatteryChargePercent: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_charge_percent"),
			"Current UPS battery charge percentage.",
			nil,
			nil,
		),

		LineVolts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "line_volts"),
			"Current AC input line voltage.",
			nil,
			nil,
		),

		LineNominalVolts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "line_nominal_volts"),
			"Nominal AC input line voltage.",
			nil,
			nil,
		),

		BatteryVolts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_volts"),
			"Current UPS battery voltage.",
			nil,
			nil,
		),

		BatteryNominalVolts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_nominal_volts"),
			"Nominal UPS battery voltage.",
			nil,
			nil,
		),

		BatteryNumberTransfersTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_number_transfers_total"),
			"Total number of transfers to UPS battery power.",
			nil,
			nil,
		),

		BatteryTimeLeftSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_time_left_seconds"),
			"Number of seconds remaining of UPS battery power.",
			nil,
			nil,
		),

		BatteryTimeOnSeconds: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_time_on_seconds"),
			"Number of seconds the UPS has been providing battery power due to an AC input line outage.",
			nil,
			nil,
		),

		BatteryCumulativeTimeOnSecondsTotal: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "battery_cumulative_time_on_seconds_total"),
			"Total number of seconds the UPS has provided battery power due to AC input line outages.",
			nil,
			nil,
		),

		LastTransferOnBattery: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "apcupsd_last_transfer_on_battery"),
			"Time of last transfer to battery since apcupsd startup.",
			nil,
			nil,
		),

		LastTransferOffBattery: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "apcupsd_last_transfer_off_battery"),
			"Time of last transfer from battery since apcupsd startup.",
			nil,
			nil,
		),

		LastSelftest: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "apcupsd_last_selftest"),
			"Time of last selftest since apcupsd startup.",
			nil,
			nil,
		),

		NominalPowerWatts: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "apcupsd_nominal_power_watts"),
			"Nominal power output in watts.",
			nil,
			nil,
		),

		UPSStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "ups_status"),
			"UPS Status values.",
			[]string{"status"},
			nil,
		),

		UPSInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "ups_info"),
			"Hostname, UPS Model and Name.",
			infoLabels,
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

	infoLabels := []string{
		s.Hostname,
		s.UPSName,
		s.Model,
	}

	upsStatus := []string{
		"CAL",
		"TRIM",
		"BOOST",
		"ONLINE",
		"ONBATT",
		"OVERLOAD",
		"LOWBATT",
		"REPLACEBATT",
		"NOBATT",
		"SLAVE",
		"SLAVEDOWN",
		"COMMLOST",
		"SHUTTING DOWN",
	}

	ch <- prometheus.MustNewConstMetric(
		c.UPSLoadPercent,
		prometheus.GaugeValue,
		s.LoadPercent,
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryChargePercent,
		prometheus.GaugeValue,
		s.BatteryChargePercent,
	)

	ch <- prometheus.MustNewConstMetric(
		c.LineVolts,
		prometheus.GaugeValue,
		s.LineVoltage,
	)

	ch <- prometheus.MustNewConstMetric(
		c.LineNominalVolts,
		prometheus.GaugeValue,
		s.NominalInputVoltage,
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryVolts,
		prometheus.GaugeValue,
		s.BatteryVoltage,
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryNominalVolts,
		prometheus.GaugeValue,
		s.NominalBatteryVoltage,
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryNumberTransfersTotal,
		prometheus.CounterValue,
		float64(s.NumberTransfers),
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryTimeLeftSeconds,
		prometheus.GaugeValue,
		s.TimeLeft.Seconds(),
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryTimeOnSeconds,
		prometheus.GaugeValue,
		s.TimeOnBattery.Seconds(),
	)

	ch <- prometheus.MustNewConstMetric(
		c.BatteryCumulativeTimeOnSecondsTotal,
		prometheus.CounterValue,
		s.CumulativeTimeOnBattery.Seconds(),
	)

	collectTimestamp(
		ch,
		c.LastTransferOnBattery,
		s.XOnBattery,
	)

	collectTimestamp(
		ch,
		c.LastTransferOffBattery,
		s.XOffBattery,
	)

	collectTimestamp(
		ch,
		c.LastSelftest,
		s.LastSelftest,
	)

	ch <- prometheus.MustNewConstMetric(
		c.NominalPowerWatts,
		prometheus.GaugeValue,
		float64(s.NominalPower),
	)

	for _, status := range upsStatus {
		value := float64(0)
		if strings.Contains(s.Status, status) {
			value = float64(1)
		}
		ch <- prometheus.MustNewConstMetric(
			c.UPSStatus,
			prometheus.GaugeValue,
			value,
			status,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		c.UPSInfo,
		prometheus.GaugeValue,
		float64(1),
		infoLabels...,
	)

	return nil, nil
}

// collectTimestamp collects timestamp metrics.
// Timestamps that are zero (time.IsZero() == true) are ignored, as such a timestamp indicates
// 'information not available', which is best expressed in Prometheus by not having the metric at all.
func collectTimestamp(ch chan<- prometheus.Metric, desc *prometheus.Desc, time time.Time, labelValues ...string) {
	if time.IsZero() {
		return
	}

	ch <- prometheus.MustNewConstMetric(
		desc,
		prometheus.GaugeValue,
		float64(time.Unix()),
		labelValues...,
	)
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
		c.UPSStatus,
		c.UPSInfo,
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
