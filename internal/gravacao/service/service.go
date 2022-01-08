package service

import (
	"bytes"
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/filipeandrade6/vigia-go/internal/core/camera"
	"github.com/filipeandrade6/vigia-go/internal/core/processo"
	"github.com/filipeandrade6/vigia-go/internal/core/registro"
	"github.com/filipeandrade6/vigia-go/internal/core/veiculo"
	"github.com/filipeandrade6/vigia-go/internal/sys/operrors"
	"github.com/nats-io/nats.go"

	"go.uber.org/zap"
)

type Service struct {
	log  *zap.SugaredLogger
	msgr *nats.Conn

	mu        *sync.RWMutex
	processos map[string]Camera // TODO era ponteiros
	retry     map[string]Camera // TODO era ponteiros
	matchlist map[string]bool

	cameraCore   camera.Core
	processoCore processo.Core
	registroCore registro.Core
	veiculoCore  veiculo.Core
}

func NewService(
	log *zap.SugaredLogger,
	msgr *nats.Conn,
	cameraCore camera.Core,
	processoCore processo.Core,
	registroCore registro.Core,
	veiculoCore veiculo.Core,
) *Service {
	return &Service{
		log:          log,
		msgr:         msgr,
		cameraCore:   cameraCore,
		processoCore: processoCore,
		registroCore: registroCore,
		veiculoCore:  veiculoCore,
	}
}

func (s *Service) Start() {
	_, err := s.msgr.Subscribe("management", s.managementH)
	if err != nil {
		fmt.Println(err)
	}

	_, err = s.msgr.Subscribe("registro", s.registroH)
	if err != nil {
		fmt.Println(err)
	}

	runtime.Goexit()
}

func (s *Service) managementH(msg *nats.Msg) {
	switch {
	case bytes.HasPrefix(msg.Data, []byte("processo.create")):

	}

	fmt.Println("veiculo recebido", string(msg.Data))
}

func (s *Service) registroH(msg *nats.Msg) {
	fmt.Println("hello")
}

func (s *Service) createAndCheckRegistro(msg *nats.Msg) {

	// Tratar msg para mandar para o banco de dados

	_, err := s.registroCore.Create(context.Background(), reg)
	if err != nil {
		if nonStoppedPrc := p.StopProcessos([]string{reg.ProcessoID}); nonStoppedPrc != nil {
			p.errChan <- operrors.OpError{Err: fmt.Errorf("could not stop processo: %s", nonStoppedPrc)}
		}
		p.errChan <- operrors.OpError{Err: fmt.Errorf("could not create registro: %w", err)}
		return
	}

	p.mu.RLock()
	_, ok := p.matchlist[reg.Placa]
	p.mu.RUnlock()
	if ok {
		p.matchChan <- reg.RegistroID
	}
}
