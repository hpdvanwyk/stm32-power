/*
Copyright (c) 2020 Hendrik van Wyk
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
this list of conditions and the following disclaimer in the documentation
and/or other materials provided with the distribution.

* Neither the name of the copyright holder nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package main

import (
	"fmt"
	"math"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/hpdvanwyk/stm32-power/blob/master/pc_client/pb"

	"github.com/prometheus/client_golang/prometheus"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

const DefaultExpiryTime = 100 * time.Second

type SensorConfig struct {
	Rename  map[string]string
	Fudge   map[string]float64
	Timeout map[string]int
	Vcal    float64
	Ical    map[int]float64
}

var defaultConfig = &SensorConfig{
	Vcal: 1,
	Ical: map[int]float64{
		0: 1,
		1: 1,
		2: 1,
	},
}

func IdString(id []byte) string {
	return net.HardwareAddr(id).String()
}

type label struct {
	key   string
	value string
}

type metric struct {
	totalValues float64
	count       int
	lables      []*label
	opts        *prometheus.GaugeOpts
	lastUpdated time.Time
	expiryTime  time.Duration
}

func (s *SensorExporter) IdString(id []byte) string {
	idString := IdString(id)
	rename, hasRename := s.config.Rename[idString]
	if hasRename {
		return rename
	}
	return idString
}

func (s *SensorExporter) fudgeTemperature(id []byte) float64 {
	idString := IdString(id)
	fudge, hasFudge := s.config.Fudge[idString]
	if hasFudge {
		return fudge
	}
	return 0
}

func (s *SensorExporter) expiryTime(id []byte) time.Duration {
	idString := IdString(id)
	timeout, hasTimeout := s.config.Timeout[idString]
	if hasTimeout {
		return time.Duration(timeout) * time.Second
	}
	return DefaultExpiryTime
}

func (se *SensorExporter) updateMetric(value float64,
	opts *prometheus.GaugeOpts,
	expiryTime time.Duration,
	labels ...*label) {
	se.Lock()
	defer se.Unlock()
	k := opts.Name + "~"
	for i := range labels {
		k += labels[i].key + "~" + labels[i].value + "~"
	}
	t, ok := se.metrics[k]
	if !ok {
		t = &metric{
			lables:     labels,
			opts:       opts,
			expiryTime: expiryTime,
		}
		se.metrics[k] = t
	}
	t.totalValues += value
	t.count++
	t.lastUpdated = time.Now()
}

func (se *SensorExporter) gc() {
	se.Lock()
	defer se.Unlock()
	for i := range se.metrics {
		if time.Since(se.metrics[i].lastUpdated) > se.metrics[i].expiryTime {
			delete(se.metrics, i)
		}
	}
}

type energyTrack struct {
	energyUsed  float64
	lastUpdated time.Time
}

type SensorExporter struct {
	sync.RWMutex
	MsgChan chan *pb.PowerMessage
	close   chan struct{}
	metrics map[string]*metric
	config  *SensorConfig
	energy  map[int]*energyTrack
	msg     *pb.PowerMessage
}

var VoltageOpts = prometheus.GaugeOpts{
	Name: "power_sensor_voltage_rms_v",
	Help: "RMS voltage of mains.",
}

var CurrentOpts = prometheus.GaugeOpts{
	Name: "power_sensor_current_rms_a",
	Help: "RMS current",
}

var RealPowerOpts = prometheus.GaugeOpts{
	Name: "power_sensor_real_power_w",
	Help: "Real power",
}

var ApparentPowerOpts = prometheus.GaugeOpts{
	Name: "power_sensor_apparent_power_va",
	Help: "Apparent power",
}

var PowerFactorOpts = prometheus.GaugeOpts{
	Name: "power_sensor_power_factor",
	Help: "Power factor",
}

func NewSensorExporter(config *SensorConfig) *SensorExporter {
	s := &SensorExporter{
		MsgChan: make(chan *pb.PowerMessage, 2),
		close:   make(chan struct{}),
		metrics: make(map[string]*metric),
		energy:  make(map[int]*energyTrack),
		config:  config,
	}
	return s
}

func (s *SensorExporter) Describe(chan<- *prometheus.Desc) {
}

func (se *SensorExporter) Collect(metricChan chan<- prometheus.Metric) {
	se.Lock()
	defer se.Unlock()
	for i, m := range se.metrics {
		opts := *se.metrics[i].opts
		labels := make(map[string]string)
		for j := range m.lables {
			labels[m.lables[j].key] = m.lables[j].value
		}
		opts.ConstLabels = labels
		g := prometheus.NewGauge(opts)
		// Yes I know this technically breaks the prometheus contract
		// Maybe I should consider a push based tsdb for these kinds of
		// measurements?
		g.Set(m.totalValues / float64(m.count))
		m.totalValues = 0
		m.count = 0
		metricChan <- g
	}
}

func (s *SensorExporter) exportPowerReading(pm *pb.PowerMessage) {
	expiryTime := 30 * time.Second
	vCal := s.config.Vcal
	s.updateMetric(
		float64(pm.VoltageRms)/vCal,
		&VoltageOpts,
		expiryTime,
	)
	fmt.Printf("Voltage %v\n", float64(pm.VoltageRms)/vCal)
	for i, p := range pm.Powers {
		iCal := s.config.Ical[i]
		s.updateMetric(
			float64(p.RealPower)/(iCal*vCal),
			&RealPowerOpts,
			expiryTime,
			&label{"sensor", strconv.Itoa(i)},
		)
		s.updateMetric(
			float64(p.ApparentPower)/(iCal*vCal),
			&ApparentPowerOpts,
			expiryTime,
			&label{"sensor", strconv.Itoa(i)},
		)
		s.updateMetric(
			float64(p.CurrentRms)/(iCal),
			&CurrentOpts,
			expiryTime,
			&label{"sensor", strconv.Itoa(i)},
		)
		s.updateMetric(
			float64(p.PowerFactor),
			&PowerFactorOpts,
			expiryTime,
			&label{"sensor", strconv.Itoa(i)},
		)
		e, exists := s.energy[i]
		if !exists {
			s.energy[i] = &energyTrack{
				lastUpdated: time.Now(),
			}
			e = s.energy[i]
		}
		now := time.Now()
		timePassed := now.Sub(e.lastUpdated)
		e.lastUpdated = now
		e.energyUsed += (timePassed.Hours() * float64(p.RealPower) / (iCal * vCal)) / 1000

		fmt.Printf("%v RealPower %v\n", i, float64(p.RealPower)/(iCal*vCal))
		fmt.Printf("%v ApparentPower %v\n", i, float64(p.ApparentPower)/(iCal*vCal))
		fmt.Printf("%v Current %v\n", i, float64(p.CurrentRms)/(iCal))
		fmt.Printf("%v Power factor %v\n", i, p.PowerFactor)
		fmt.Printf("%v Enegry %v kW h\n", i, e.energyUsed)
	}
}

func (s *SensorExporter) Run() {
	cleanup := time.NewTicker(1 * time.Second)
	for {
		select {
		case msg := <-s.MsgChan:
			s.exportPowerReading(msg)
			s.Lock()
			s.msg = msg
			s.Unlock()
		case <-cleanup.C:
			s.gc()
		case <-s.close:
			cleanup.Stop()
			return
		}
	}
}

func (s *SensorExporter) HandleChart(source int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.RLock()
		defer s.RUnlock()
		if s.msg == nil {
			return
		}
		iCal := s.config.Ical[source]
		vCal := s.config.Vcal
		current := make(plotter.XYs, len(s.msg.Powers[source].Current))
		maxCurrent := 0.0
		fmt.Printf("DC i %v\n", s.msg.Powers[source].DC)
		fmt.Printf("DC v %v\n", s.msg.DC)
		for i := range current {
			current[i].X = float64(i)
			current[i].Y = (float64(s.msg.Powers[source].Current[i]) - float64(s.msg.Powers[source].DC)) / iCal
			if math.Abs(current[i].Y) > maxCurrent {
				maxCurrent = math.Abs(current[i].Y)
			}
		}

		v := make(plotter.XYs, len(s.msg.Voltage))
		for i := range v {
			v[i].X = float64(i)
			v[i].Y = (((float64(s.msg.Voltage[i]) - float64(s.msg.DC)) / vCal) / 350) * maxCurrent
		}

		p, err := plot.New()
		if err != nil {
			panic(err)
		}

		p.Title.Text = "Current"
		p.X.Label.Text = "Sample"
		p.Y.Label.Text = "A"
		var xTics plot.ConstantTicks = make([]plot.Tick, 0)
		for i := 0; i < len(current); i += 20 {
			xTics = append(xTics, plot.Tick{
				Value: float64(i),
				Label: strconv.Itoa(i),
			})
		}
		p.X.Tick.Marker = xTics
		grid := plotter.NewGrid()
		p.Add(grid)
		err = plotutil.AddLines(
			p,
			"Current", current,
			"Voltage", v,
		)

		if err != nil {
			panic(err)
		}

		wt, err := p.WriterTo(10*vg.Inch, 6*vg.Inch, "svg")
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		w.Header().Set("Content-Type", "image/svg+xml")
		wt.WriteTo(w)
	}
}

func (s *SensorExporter) HandleChartVoltage(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	defer s.RUnlock()
	if s.msg == nil {
		return
	}
	vCal := s.config.Vcal

	v := make(plotter.XYs, len(s.msg.Voltage))
	for i := range v {
		v[i].X = float64(i)
		v[i].Y = ((float64(s.msg.Voltage[i]) - float64(s.msg.DC)) / vCal)
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Voltage"
	p.X.Label.Text = "Sample"
	p.Y.Label.Text = "Measurement"
	var xTics plot.ConstantTicks = make([]plot.Tick, 0)
	for i := 0; i < len(v); i += 20 {
		xTics = append(xTics, plot.Tick{
			Value: float64(i),
			Label: strconv.Itoa(i),
		})
	}
	var yTics plot.ConstantTicks = make([]plot.Tick, 0)
	for i := -360; i < 360; i += 40 {
		yTics = append(yTics, plot.Tick{
			Value: float64(i),
			Label: strconv.Itoa(i),
		})
	}
	p.X.Tick.Marker = xTics
	p.Y.Tick.Marker = yTics
	p.Add(plotter.NewGrid())
	err = plotutil.AddLines(
		p,
		"Voltage", v,
	)

	if err != nil {
		panic(err)
	}

	wt, err := p.WriterTo(10*vg.Inch, 6*vg.Inch, "svg")
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	wt.WriteTo(w)
}
