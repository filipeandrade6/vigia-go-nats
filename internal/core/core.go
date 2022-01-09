package core

import (
	"expvar"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/filipeandrade6/vigia-go/internal/core/camera"
	"github.com/filipeandrade6/vigia-go/internal/core/config"
	"github.com/filipeandrade6/vigia-go/internal/core/processo"
	"github.com/filipeandrade6/vigia-go/internal/core/registro"
	"github.com/filipeandrade6/vigia-go/internal/core/servidorgravacao"
	"github.com/filipeandrade6/vigia-go/internal/core/veiculo"
	"github.com/filipeandrade6/vigia-go/internal/sys/database"
	"github.com/filipeandrade6/vigia-go/internal/sys/messaging"
	"github.com/nats-io/nats.go"

	"github.com/filipeandrade6/vigia-go/internal/gravacao/service"

	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

var build = "develop"

func Run(log *zap.SugaredLogger, cfg config.Configuration) error {
	// =========================================================================
	// CPU quota

	if _, err := maxprocs.Set(); err != nil {
		log.Errorw("startup", zap.Error(err))
		os.Exit(1)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// =========================================================================
	// Show configuration

	log.Infow("startup", "config", fmt.Sprintf("%+v", cfg)) // TODO esconder senhas

	// =========================================================================
	// App starting

	expvar.NewString("build").Set(build)   // TODO expvar
	log.Infow("startup", "version", build) // TODO utilizar essa build
	defer log.Infow("shutdown complete")

	// =========================================================================
	// TODO Initialize Authentication Support

	// =========================================================================
	// Database support

	db, err := database.Connect(database.Config{
		User:         cfg.Database.User,
		Password:     cfg.Database.Password,
		Host:         cfg.Database.Host,
		Name:         cfg.Database.Name,
		MaxIDLEConns: cfg.Database.MaxIDLEConns,
		MaxOpenConns: cfg.Database.MaxOpenConns,
		DisableTLS:   cfg.Database.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Infow("shutdown", "status", "stopping database support", "host", cfg.Database.Host)
		db.Close()
	}()

	// =========================================================================
	// Messaging support

	msgrConn, err := messaging.Connect(messaging.Config{
		Name: "manager",               // TODO puxar das configs
		URL:  "nats://localhost:4222", // TODO puxar das configs
	})
	if err != nil {
		return fmt.Errorf("connecting to messaging support: %w", err)
	}

	msgr, err := nats.NewEncodedConn(msgrConn, nats.JSON_ENCODER)
	if err != nil {
		return fmt.Errorf("wrapping with encoder the messaging support: %w", err)
	}

	defer func() {
		log.Infow("shutdown", "status", "stopping messaging support", "host", "cfg.Messaging.URL") // TODO trocar

		if err := msgr.Flush(); err != nil {
			log.Errorf("flushing the messaging support: %w", err)
		}

		if err := msgr.LastError(); err != nil {
			log.Errorf("last error from messaging support: %w", err)
		}

		msgr.Close()
	}()

	// TODO tratar errors de reconexão, talvez um handler, etc
	// https://github.com/nats-io/nats.go/blob/main/examples/nats-sub/main.go

	// =========================================================================
	// Start service

	log.Infow("startup", "status", "initializing service")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	servidorGravacaoCore := servidorgravacao.NewCore(log, db)
	cameraCore := camera.NewCore(log, db)
	processoCore := processo.NewCore(log, db)
	registroCore := registro.NewCore(log, db)
	veiculoCore := veiculo.NewCore(log, db)

	svc := service.NewService(
		log,
		msgr,
		servidorGravacaoCore,
		cameraCore,
		processoCore,
		registroCore,
		veiculoCore,
	)

	go func() {
		svc.Start()
	}()

	// =========================================================================
	// Shutdown

	// select {
	// case err := <-serverErrors:
	// 	return fmt.Errorf("server error: %w", err)

	sig := <-shutdown
	log.Infow("shutdown", "status", "shutdown started", "signal", sig)
	defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)
	// TODO gracefull shutdown the service

	return nil
}
