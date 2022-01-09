package service

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/filipeandrade6/vigia-go/internal/core/camera"
	"github.com/filipeandrade6/vigia-go/internal/core/processo"
	"github.com/filipeandrade6/vigia-go/internal/core/registro"
	"github.com/filipeandrade6/vigia-go/internal/core/servidorgravacao"
	"github.com/filipeandrade6/vigia-go/internal/core/veiculo"
	"github.com/nats-io/nats.go"

	"go.uber.org/zap"
)

type Service struct {
	log  *zap.SugaredLogger
	msgr *nats.EncodedConn

	mu                 sync.RWMutex
	matchlist          map[string]bool // TODO trocar por struct?
	servidoresGravacao map[string]bool // ! necessário? TODO trocar por struct?
	processos          map[string]processo.Processo

	servidorGravacaoCore servidorgravacao.Core
	cameraCore           camera.Core
	processoCore         processo.Core
	registroCore         registro.Core
	veiculoCore          veiculo.Core
}

func NewService(
	log *zap.SugaredLogger,
	msgr *nats.EncodedConn,
	servidorGravacaoCore servidorgravacao.Core,
	cameraCore camera.Core,
	processoCore processo.Core,
	registroCore registro.Core,
	veiculoCore veiculo.Core,

) *Service {
	return &Service{
		log:                  log,
		msgr:                 msgr,
		matchlist:            make(map[string]bool), // TODO trocar por struct?
		servidoresGravacao:   make(map[string]bool), // ! necessário? TODO trocar por struct?
		processos:            make(map[string]processo.Processo),
		servidorGravacaoCore: servidorGravacaoCore,
		cameraCore:           cameraCore,
		processoCore:         processoCore,
		registroCore:         registroCore,
		veiculoCore:          veiculoCore,
	}
}

func (s *Service) Start() {

	// External
	_, err := s.msgr.Subscribe("management.>", s.managementHandler)
	if err != nil {
		fmt.Println(err)
	}

	// Internal
	_, err = s.msgr.Subscribe("descoberta", s.descobertaHandler)
	if err != nil {
		fmt.Println(err)
	}

	_, err = s.msgr.Subscribe("registro.>", s.registroHandler)
	if err != nil {
		fmt.Println(err)
	}

	_, err = s.msgr.Subscribe("erro.>", s.erroHandler)
	if err != nil {
		fmt.Println(err)
	}

	runtime.Goexit()
}
