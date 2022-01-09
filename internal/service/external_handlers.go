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
	switch msg.Subject {
	case "management.processo.iniciar":
		fmt.Println("recebido processo.create")
		if err := s.iniciarProcesso(string(msg.Data)); err != nil {
			s.log.Errorf("handling management.processo.iniciar", err)
		}

	case "management.processo.parar":
		fmt.Println("recebido processo.parar")
		if err := s.pararProcesso(string(msg.Data)); err != nil {
			s.log.Errorf("handling management.processo.iniciar", err)
		}

	case "management.processo.listar":
		fmt.Println("recebido processo.listar")
		prcs, err := s.listarProcesso()
		if err != nil {
			s.log.Errorf("handling management.processo.iniciar", err)
		}
		if err := s.msgr.Publish(msg.Reply, prcs); err != nil {
			s.log.Errorf("replying management.processo.iniciar", err)
		}

	case "management.armazenamento.atualizar":
		fmt.Println("recebido armazenamento.atualizar")
		if err := s.atualizarArmazenamento(); err != nil {
			s.log.Errorf("handling management.processo.iniciar", err)
		}

	case "management.match.atualizar":
		fmt.Println("atualizar lista de match")
		if err := s.atualizarMatch(); err != nil {
			s.log.Errorf("handling management.processo.iniciar", err)
		}
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
		return fmt.Errorf("querying processo in db: %w", err)
	}

	cam, err := s.cameraCore.QueryByID(context.Background(), prc.CameraID)
	if err != nil {
		return fmt.Errorf("querying camera in db: %w", err)
	}

	camByte, err := json.Marshal(cam)
	if err != nil {
		return fmt.Errorf("marshalling camera data: %w", err)
	}

	// TODO PRECISO PASSAR TANTO O PROCESSOID QUANDO CAMERA
	err = s.msgr.Publish(prc.ServidorGravacaoID+".iniciar", camByte)
	if err != nil {
		return fmt.Errorf("sending message to processo: %w", err)
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
		return fmt.Errorf("sending message to processo: %w", err)
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
			return nil, fmt.Errorf("sending request message: %w", err)
		}
		prcResponse = append(prcResponse, prcs...)
	}

	prcResponseByte, err := json.Marshal(prcResponse)
	if err != nil {
		return nil, fmt.Errorf("marshalling processo data: %w", err)
	}

	return prcResponseByte, nil
}

func (s *Service) atualizarArmazenamento() error {
	fmt.Println("implementar a atualização de armazenamento")
	return nil
}

func (s *Service) atualizarMatch() error {
	veiculos, err := s.veiculoCore.QueryAll(context.Background())
	if err != nil {
		return fmt.Errorf("querying veiculos in db: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.matchlist = make(map[string]bool)
	for _, v := range veiculos {
		s.matchlist[v.Placa] = true
	}

	return nil
}
