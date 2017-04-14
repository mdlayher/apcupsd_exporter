package apcupsdexporter

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/mdlayher/apcupsd"
)

func TestUPSCollector(t *testing.T) {
	var tests = []struct {
		desc    string
		ss      *testStatusSource
		matches []*regexp.Regexp
	}{
		{
			desc: "empty",
			ss: &testStatusSource{
				s: &apcupsd.Status{},
			},
			matches: []*regexp.Regexp{
				regexp.MustCompile(`apcupsd_battery_charge_percent{hostname="",model="",ups_name=""} 0`),
			},
		},
		{
			desc: "full",
			ss: &testStatusSource{
				s: &apcupsd.Status{
					Hostname: "a",
					Model:    "b",
					UPSName:  "c",

					BatteryChargePercent:    100.0,
					CumulativeTimeOnBattery: 30 * time.Second,
					NominalBatteryVoltage:   12.0,
					TimeLeft:                2 * time.Minute,
					TimeOnBattery:           10 * time.Second,
					BatteryVoltage:          13.2,
					NominalInputVoltage:     120.0,
					LineVoltage:             121.1,
					LoadPercent:             16.0,
					NumberTransfers:         1,
					XOnBattery:              time.Unix(100001, 0),
					XOffBattery:             time.Unix(100002, 0),
					LastSelftest:            time.Unix(100003, 0),
					NominalPower:            50.0,
				},
			},
			matches: []*regexp.Regexp{
				regexp.MustCompile(`apcupsd_battery_charge_percent{hostname="a",model="b",ups_name="c"} 100`),
				regexp.MustCompile(`apcupsd_battery_cumulative_time_on_seconds_total{hostname="a",model="b",ups_name="c"} 30`),
				regexp.MustCompile(`apcupsd_battery_nominal_volts{hostname="a",model="b",ups_name="c"} 12`),
				regexp.MustCompile(`apcupsd_battery_time_left_seconds{hostname="a",model="b",ups_name="c"} 120`),
				regexp.MustCompile(`apcupsd_battery_time_on_seconds{hostname="a",model="b",ups_name="c"} 10`),
				regexp.MustCompile(`apcupsd_battery_volts{hostname="a",model="b",ups_name="c"} 13.2`),
				regexp.MustCompile(`apcupsd_battery_number_transfers_total{hostname="a",model="b",ups_name="c"} 1`),
				regexp.MustCompile(`apcupsd_line_nominal_volts{hostname="a",model="b",ups_name="c"} 120`),
				regexp.MustCompile(`apcupsd_line_volts{hostname="a",model="b",ups_name="c"} 121.1`),
				regexp.MustCompile(`apcupsd_ups_load_percent{hostname="a",model="b",ups_name="c"} 16`),
				regexp.MustCompile(`apcupsd_last_transfer_on_battery{hostname="a",model="b",ups_name="c"} 100001`),
				regexp.MustCompile(`apcupsd_last_transfer_off_battery{hostname="a",model="b",ups_name="c"} 100002`),
				regexp.MustCompile(`apcupsd_last_selftest{hostname="a",model="b",ups_name="c"} 100003`),
				regexp.MustCompile(`apcupsd_nominal_power_watts{hostname="a",model="b",ups_name="c"} 50`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			out := testCollector(t, NewUPSCollector(tt.ss))

			for _, m := range tt.matches {
				name := metricName(t, m.String())

				t.Run(name, func(t *testing.T) {
					if !m.Match(out) {
						t.Fatalf("\toutput failed to match regex (regexp: %v)", m)
					}
				})
			}
		})
	}
}

// TestZeroTimesAreIgnored tests that times with a zero value (time.IsZero() == true)
// are not collected.
func TestZeroTimesAreIgnored(t *testing.T) {
	ss := &testStatusSource{
		s: &apcupsd.Status{
			XOnBattery:  time.Unix(123456, 0),
			XOffBattery: time.Time{},
		},
	}
	out := testCollector(t, NewUPSCollector(ss))
	// Test that in general timestamps are collected.
	if !regexp.MustCompile(`apcupsd_last_transfer_on_battery.* 123456`).Match(out) {
		t.Error("non-zero timestamp is not reported properly")
	}
	// Test that zero timestamps, however, are ignored.
	if regexp.MustCompile(`apcupsd_last_transfer_off_battery`).Match(out) {
		t.Error("Zero time is reported")
	}
}

func metricName(t *testing.T, metric string) string {
	ss := strings.Split(metric, " ")
	if len(ss) != 2 {
		t.Fatalf("malformed metric: %v", metric)
	}

	if !strings.Contains(ss[0], "{") {
		return ss[0]
	}

	ss = strings.Split(ss[0], "{")
	if len(ss) != 2 {
		t.Fatalf("malformed metric: %v", metric)
	}

	return ss[0]
}

var _ StatusSource = &testStatusSource{}

type testStatusSource struct {
	s *apcupsd.Status
}

func (ss *testStatusSource) Status() (*apcupsd.Status, error) {
	return ss.s, nil
}
