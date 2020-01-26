package mh_z14a

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/common/log"

	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/tarm/serial"
)

var (
	command = []byte{0xff, 0x01, 0x86, 0x00, 0x00, 0x00, 0x00, 0x00, 0x79}
)

type MHZ14A struct {
	parent   context.Context
	cfg      *serial.Config
	interval time.Duration
	ppm      prometheus.Gauge
}

func NewMHZ14A(name string, baud int, ctx context.Context, interval time.Duration) *MHZ14A {
	_cfg := &serial.Config{
		Name: name,
		Baud: baud,
	}
	return &MHZ14A{
		parent:   ctx,
		cfg:      _cfg,
		interval: interval,
		ppm: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "room",
				Name:      "co2",
			},
		),
	}
}

func (m *MHZ14A) Watch() error {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()
	if err := m.doWatch(); err != nil {
		return err
	}
	for {
		select {
		case <-m.parent.Done():
			log.Info(fmt.Sprintf("Shutting down."))
			return nil
		case <-ticker.C:
			if err := m.doWatch(); err != nil {
				return err
			}
		}
	}
}

func (m *MHZ14A) doWatch() error {
	_ppm, err := m.get()
	log.Info(_ppm, err)
	if err != nil {
		return err
	}
	m.ppm.Set(float64(_ppm))
	return nil
}

func (m *MHZ14A) get() (uint16, error) {
	_port, err := serial.OpenPort(m.cfg)
	if err != nil {
		log.Fatal(err)
	}
	_, err = _port.Write(command)
	if err != nil {
		return 0, err
	}

	buf := make([]byte, 9)
	_, err = _port.Read(buf)
	if err != nil {
		return 0, err
	}
	if err := _port.Close(); err != nil {
		return 0, err
	}
	return (uint16(buf[2]) << 8) | uint16(buf[3]), nil
}
