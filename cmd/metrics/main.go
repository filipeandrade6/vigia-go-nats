package main

import (
	"expvar"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/filipeandrade6/vigia-go/cmd/metrics/collector"
	"github.com/filipeandrade6/vigia-go/cmd/metrics/publisher"
	expvarsrv "github.com/filipeandrade6/vigia-go/cmd/metrics/publisher/expvar"
	"github.com/filipeandrade6/vigia-go/internal/sys/config"
	"github.com/filipeandrade6/vigia-go/internal/sys/logger"
	"github.com/spf13/viper"

	// "github.com/ardanlabs/service/app/services/metrics/collector"
	// "github.com/ardanlabs/service/app/services/metrics/publisher"
	// expvarsrv "github.com/ardanlabs/service/app/services/metrics/publisher/expvar"
	// "github.com/ardanlabs/service/foundation/logger"

	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

// build is the git version of this program. It is set using build flags in the makefile.
// var build = "develop"

func main() {

	// Construct the application logger.
	log, err := logger.New("METRICS")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// =========================================================================
	// GOMAXPROCS

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// =========================================================================
	// App Starting

	viper.AutomaticEnv()
	log.Infow("startup", "config", config.PrettyPrintConfig())

	// =========================================================================
	// Start Debug Service

	log.Infow("startup", "status", "debug router started", "host", viper.GetString("VIGIA_MET_WEB_DEBUGHOST"))

	// The Debug function returns a mux to listen and serve on for all the debug
	// related endpoints. This include the standard library endpoints.

	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	// Start the service listening for debug requests.
	// Not concerned with shutting this down with load shedding.
	go func() {
		if err := http.ListenAndServe(viper.GetString("VIGIA_MET_WEB_DEBUGHOST"), mux); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host", viper.GetString("VIGIA_MET_WEB_DEBUGHOST"), "ERROR", err)
		}
	}()

	// =========================================================================
	// Start expvar Service

	exp := expvarsrv.New(
		log,
		viper.GetString("VIGIA_MET_EXPVAR_HOST"),
		viper.GetString("VIGIA_MET_EXPVAR_ROUTE"),
		time.Duration(viper.GetInt("VIGIA_MET_EXPVAR_READTIMEOUT")*int(time.Second)),
		time.Duration(viper.GetInt("VIGIA_MET_EXPVAR_WRITETIMEOUT")*int(time.Second)),
		time.Duration(viper.GetInt("VIGIA_MET_EXPVAR_IDLETIMEOUT")*int(time.Second)),
	)
	defer exp.Stop(time.Duration(viper.GetInt("VIGIA_MET_EXPVAR_SHUTDOWNTIMEOUT") * int(time.Second)))

	// =========================================================================
	// Start collectors and publishers

	// Initialize to allow for the collection of metrics.
	collector, err := collector.New(viper.GetString("VIGIA_MET_COLLECT_FROM"))
	if err != nil {
		return fmt.Errorf("starting collector: %w", err)
	}

	// Create a stdout publisher.
	// TODO: Respect the cfg.publish.to config option.
	stdout := publisher.NewStdout(log)

	// Start the publisher to collect/publish metrics.
	publish, err := publisher.New(
		log,
		collector,
		time.Duration(viper.GetInt("VIGIA_MET_PUBLISH_INTERVAL")*int(time.Second)),
		exp.Publish,
		stdout.Publish,
	)
	if err != nil {
		return fmt.Errorf("starting publisher: %w", err)
	}
	defer publish.Stop()

	// =========================================================================
	// Shutdown

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	log.Infow("shutdown", "status", "shutdown started")
	defer log.Infow("shutdown", "status", "shutdown complete")

	return nil
}
