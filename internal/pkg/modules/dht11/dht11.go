package dht11

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/common/log"

	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/nandehu0323/go-dht"
	"github.com/prometheus/client_golang/prometheus"
)

type DHT11 struct {
	parent   context.Context
	gpio     int
	interval time.Duration
	temp     prometheus.Gauge
	humid    prometheus.Gauge
}

func NewDHT11(pin int, ctx context.Context, interval time.Duration) *DHT11 {
	return &DHT11{
		parent:   ctx,
		gpio:     pin,
		interval: interval,
		temp: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "room",
				Name:      "temperature",
			},
		),
		humid: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "room",
				Name:      "humidity",
			},
		),
	}
}

func (d *DHT11) Watch() error {
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()
	if err := d.doWatch(); err != nil {
		return err
	}
	for {
		select {
		case <-d.parent.Done():
			log.Info(fmt.Sprintf("Shutting down."))
			return nil
		case <-ticker.C:
			if err := d.doWatch(); err != nil {
				return err
			}
		}
	}
}

func (d *DHT11) doWatch() error {
	_temp, _humid, err := d.get()
	if err != nil {
		return err
	}
	d.temp.Set(float64(_temp))
	d.humid.Set(float64(_humid))
	return nil
}

func (d *DHT11) get() (float32, float32, error) {
	_temp, _humid, _, err := dht.ReadDHTxxWithRetry(dht.DHT11, d.gpio, false, 10)
	if err != nil {
		return 0, 0, err
	}
	return _temp, _humid, nil
}
