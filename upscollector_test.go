package apcupsdexporter

import (
	"regexp"
	"testing"
	"time"

	"github.com/mdlayher/apcupsd"
)

func TestUPSCollector(t *testing.T) {
	tests := []struct {
		desc    string
		ss      *testStatusSource
		matches []*regexp.Regexp
	}{
		{
			desc: "empty",
			ss: &testStatusSource{
				s: &apcupsd.Status{},
			},
		},
		{
			desc: "full",
			ss: &testStatusSource{
				s: &apcupsd.Status{
					Hostname: "foo",
					Model:    "APC UPS",
					UPSName:  "bar",
					Status:   "ONLINE",

					BatteryChargePercent:    100.0,
					CumulativeTimeOnBattery: 30 * time.Second,
					NominalBatteryVoltage:   12.0,
					TimeLeft:                2 * time.Minute,
					TimeOnBattery:           10 * time.Second,
					BatteryVoltage:          13.2,
					NominalInputVoltage:     120.0,
					LineVoltage:             121.1,
					OutputVoltage:           120.9,
					OutputAmps:              1.06,
					LoadPercent:             16.0,
					NumberTransfers:         1,
					XOnBattery:              time.Unix(100001, 0),
					XOffBattery:             time.Unix(100002, 0),
					LastSelftest:            time.Unix(100003, 0),
					NominalPower:            50.0,
					InternalTemp:            26.4,
				},
			},
			matches: []*regexp.Regexp{
				regexp.MustCompile(`apcupsd_battery_charge_percent{ups="bar"} 100`),
				regexp.MustCompile(`apcupsd_battery_cumulative_time_on_seconds_total{ups="bar"} 30`),
				regexp.MustCompile(`apcupsd_battery_nominal_volts{ups="bar"} 12`),
				regexp.MustCompile(`apcupsd_battery_time_left_seconds{ups="bar"} 120`),
				regexp.MustCompile(`apcupsd_battery_time_on_seconds{ups="bar"} 10`),
				regexp.MustCompile(`apcupsd_battery_volts{ups="bar"} 13.2`),
				regexp.MustCompile(`apcupsd_battery_number_transfers_total{ups="bar"} 1`),
				regexp.MustCompile(`apcupsd_info{hostname="foo",model="APC UPS",status="ONLINE",ups="bar"} 1`),
				regexp.MustCompile(`apcupsd_line_nominal_volts{ups="bar"} 120`),
				regexp.MustCompile(`apcupsd_line_volts{ups="bar"} 121.1`),
				regexp.MustCompile(`apcupsd_output_volts{ups="bar"} 120.9`),
				regexp.MustCompile(`apcupsd_output_amps{ups="bar"} 1.06`),
				regexp.MustCompile(`apcupsd_ups_load_percent{ups="bar"} 16`),
				regexp.MustCompile(`apcupsd_last_transfer_on_battery_time_seconds{ups="bar"} 100001`),
				regexp.MustCompile(`apcupsd_last_transfer_off_battery_time_seconds{ups="bar"} 100002`),
				regexp.MustCompile(`apcupsd_last_selftest_time_seconds{ups="bar"} 100003`),
				regexp.MustCompile(`apcupsd_nominal_power_watts{ups="bar"} 50`),
				regexp.MustCompile(`apcupsd_internal_temperature_celsius{ups="bar"} 26.4`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			out := testCollector(t, NewUPSCollector(tt.ss))

			for _, m := range tt.matches {
				if !m.Match(out) {
					t.Fatalf("output failed to match regex (regexp: %v)", m)
				}
			}
		})
	}
}

var _ StatusSource = &testStatusSource{}

type testStatusSource struct {
	s *apcupsd.Status
}

func (ss *testStatusSource) Status() (*apcupsd.Status, error) {
	return ss.s, nil
}
