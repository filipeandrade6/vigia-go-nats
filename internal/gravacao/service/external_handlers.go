package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/filipeandrade6/vigia-go/internal/core/processo"
	"github.com/nats-io/nats.go"
)

func (s *Service) managementHandler(msg *nats.Msg) {
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

func (s *Service) iniciarProcesso(prcID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.processos[prcID]
	if ok {
		return errors.New("alerady executing")
	}

	prc, err := s.processoCore.QueryByID(context.Background(), prcID)
	if err != nil {
		fmt.Println(err) // ! tratar
	}

	cam, err := s.cameraCore.QueryByID(context.Background(), prc.CameraID)
	if err != nil {
		fmt.Println(err) // ! tratar
	}

	camByte, err := json.Marshal(cam)
	if err != nil {
		fmt.Println(err) // ! tratar
	}

	// TODO PRECISO PASSAR TANTO O PROCESSOID QUANDO CAMERA
	err = s.msgr.Publish(prc.ServidorGravacaoID+".iniciar", camByte)
	if err != nil {
		fmt.Println(err) // ! tratar
	}

	s.processos[prcID] = prc

	return nil
}

func (s *Service) pararProcesso(prcID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	prc, ok := s.processos[prcID]
	if !ok {
		return errors.New("processo is not running")
	}

	if err := s.msgr.Publish(prc.ServidorGravacaoID+".parar", []byte(prc.ProcessoID)); err != nil {
		fmt.Println(err) // ! arrumar
	}

	return nil
}

func (s *Service) listarProcesso() ([]byte, error) {
	prcResponse := []processo.Processo{}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for sv := range s.servidoresGravacao {
		prcs := []processo.Processo{}
		if err := s.msgr.Request(sv+".listar", nil, &prcs, 500*time.Millisecond); err != nil {
			return nil, err // !tratar
		}
		prcResponse = append(prcResponse, prcs...)
	}

	prcResponseByte, err := json.Marshal(prcResponse)
	if err != nil {
		return nil, err // !tratar
	}

	return prcResponseByte, nil
}

func (s *Service) atualizarArmazenamento() {
	fmt.Println("implementar a atualização de armazenamento")
}

func (s *Service) atualizarMatch() error {
	veiculos, err := s.veiculoCore.QueryAll(context.Background())
	if err != nil {
		fmt.Println(err) // ! arrumar
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.matchlist = make(map[string]bool)
	for _, v := range veiculos {
		s.matchlist[v.Placa] = true
	}

	return nil
}
