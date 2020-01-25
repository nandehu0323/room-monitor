package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nandehu0323/room-monitor/application"

	mh_z14a "github.com/nandehu0323/room-monitor/internal/pkg/modules/mh-z14a"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"

	"github.com/nandehu0323/room-monitor/internal/pkg/modules/dht11"
)

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer close(sigCh)
		for sig := range sigCh {
			log.Info(fmt.Sprintf("SIGNAL %d received. Then canceling jobs.", sig))
			cancel()
			break
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:    ":19100",
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Info(fmt.Sprintf("Exporter running on: %s\n", "http://localhost:19100/metrics"))

	monitor := application.NewMonitor(ctx)
	monitor.Register(dht11.NewDHT11(4, ctx, 5*time.Second))
	monitor.Register(mh_z14a.NewMHZ14A("/dev/ttyAMA0", 9600, ctx, 5*time.Second))

	if err := monitor.Run(); err != nil {
		log.Fatal(err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Info("Server Exited Properly")
}
