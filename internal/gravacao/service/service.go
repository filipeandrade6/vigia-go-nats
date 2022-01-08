package service

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"

	"github.com/filipeandrade6/vigia-go/internal/core/camera"
	"github.com/filipeandrade6/vigia-go/internal/core/processo"
	"github.com/filipeandrade6/vigia-go/internal/core/registro"
	"github.com/filipeandrade6/vigia-go/internal/core/veiculo"
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
	msgr *nats.EncodedConn,
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
	switch msg.Subject {
	case "processo.iniciar":
		fmt.Println("recebido processo.create") // TODO COLOCAR REPLY?
		s.iniciarProcesso(string(msg.Data))
	case "processo.parar":
		fmt.Println("recebido processo.parar") // TODO COLOCAR REPLY?
		s.pararProcesso(string(msg.Data))
	case "processo.listar":
		fmt.Println("recebido processo.listar")
		s.listarProcesso()
	}
}

func (s *Service) iniciarProcesso(camID string) {
	cam, err := s.cameraCore.QueryByID(context.Background(), camID)
	if err != nil {
		fmt.Println("ERROR") // ! arrumar
	}

	camByte, err := json.Marshal(cam)
	if err != nil {
		fmt.Println("ERROR")
	}

	err = s.msgr.Publish("camera.start", camByte)
	if err != nil {
		fmt.Println("ERROR") // ! arrumar
	}
}

func (s *Service) pararProcesso(camID string) {
	if err := s.msgr.Publish("camera.stop", []byte(camID)); err != nil {
		fmt.Println("ERROR") // ! arrumar
	}
}

func (s *Service) listarProcesso() {
	// TODO colocar reply

	if err := s.msgr.Publish("camera.listar", nil); err != nil {
		fmt.Println("ERROR") // ! arrumar
	}
	fmt.Println()
}

func (s *Service) registroH(msg *nats.Msg) {
	reg := registro.Registro{}
	err := json.Unmarshal(msg.Data, &reg)
	if err != nil {
		fmt.Println("ERROR") // ! arrumar
	}

	_, err = s.registroCore.Create(context.Background(), reg)
	if err != nil {
		fmt.Println("nao pode criar o registro, ver qual tratamento")
		return
	}

	s.mu.RLock()
	defer s.mu.Unlock()
	_, ok := s.matchlist[reg.Placa]
	if ok {
		alerta(reg.RegistroID)
	}
}

func alerta(regID string) {
	fmt.Println("encontrado - alertaID:", regID)
}
