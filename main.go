package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/vvrnv/kube-ns-cleaner/config"
	"github.com/vvrnv/kube-ns-cleaner/k8sclient"
	"github.com/vvrnv/kube-ns-cleaner/metrics"
)

// scheduler for running all functions
func runCronJobs() {
	scheduler := gocron.NewScheduler(time.Local)

	cron := config.Config.Cron
	scheduler.Cron(cron).Do(func() {
		log.Info().Msg("Scheduler started")
		k8sclient.K8sClient()
	})

	scheduler.StartBlocking()
}

func statusHandler(logger zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Service is working")
		log.Info().Msgf("GET %s 200", r.URL.Path)
	}
}

func main() {
	// set up logging
	runLogFile, _ := os.OpenFile(
		//"/opt/app/log/kube-ns-cleaner.json",
		"kube-ns-cleaner.json", // mac os debug
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	multi := zerolog.MultiLevelWriter(os.Stdout, runLogFile)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05.000"})

	log.Info().Msg(("Kube-ns-cleaner is starting"))

	// load application configuration
	if err := config.LoadConfig("/opt/app/"); err != nil {
		log.Fatal().Msgf("invalid application configuration: %s", err)
	}

	log.Info().Msg(("Config successfully loaded"))
	log.Info().Msgf("Waiting for a cron schedule. Current schedule is `%s`", config.Config.Cron)

	http.HandleFunc("/api/status", statusHandler(log.Logger))
	metrics.RecordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	addr := ":5001"
	log.Info().Msgf("Starting server at %s", addr)

	go runCronJobs()

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal().Msgf("Server stopped %s", err)
	}
}
