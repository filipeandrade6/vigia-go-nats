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
	msgr *nats.EncodedConn

	mu        sync.RWMutex
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
	_, err := s.msgr.Subscribe("management.>", s.managementH)
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
	fmt.Println(msg.Subject)

	switch msg.Subject {
	case "management.processo.iniciar":
		fmt.Println("recebido processo.create") // TODO COLOCAR REPLY?
		s.iniciarProcesso(string(msg.Data))

	case "management.processo.parar":
		fmt.Println("recebido processo.parar") // TODO COLOCAR REPLY?
		s.pararProcesso(string(msg.Data))

	case "management.processo.listar":
		fmt.Println("recebido processo.listar") // TODO COLOCAR REPLY?
		s.listarProcesso()

	case "management.armazenamento.atualizar":
		fmt.Println("recebido armazenamento.atualizar") // TODO COLOCAR REPLY?
		s.atualizarArmazenamento()

	case "management.match.atualizar":
		fmt.Println("atualizar lista de match")
		s.atualizarMatch()
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
	fmt.Println("entrou listar Processos")
}

func (s *Service) atualizarArmazenamento() {
	fmt.Println("implementar a atualização de armazenamento")
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

func (s *Service) atualizarMatch() {
	veiculos, err := s.veiculoCore.QueryAll(context.Background())
	if err != nil {
		fmt.Println("ERROR") // ! arrumar
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.matchlist = make(map[string]bool)
	for _, v := range veiculos {
		s.matchlist[v.Placa] = true
	}
}
