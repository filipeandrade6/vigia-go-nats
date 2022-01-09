package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/filipeandrade6/vigia-go/internal/core/registro"
	"github.com/filipeandrade6/vigia-go/internal/core/servidorgravacao"
	"github.com/nats-io/nats.go"
)

// * Se descoberta retornar nula é pq tem erro

func (s *Service) descobertaHandler(msg *nats.Msg) {
	// ? ver se vai funcionar isso aqui
	var err error
	defer func(err error) {
		if err != nil {
			if err = msg.Respond(nil); err != nil {
				s.log.Errorf("responding mensage: %s", err)
			}
		}
	}(err)

	var sv servidorgravacao.NewServidorGravacao

	err = json.Unmarshal(msg.Data, &sv)
	if err != nil {
		s.log.Errorf("unmarshalling mensage: %s", err)
		return
	}

	svDB, err := s.servidorGravacaoCore.Create(context.Background(), sv)
	if err != nil {
		s.log.Errorf("creating servidor gravacao in database: %s", err)
		return
	}

	err = msg.Respond([]byte(svDB.ServidorGravacaoID))
	if err != nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.servidoresGravacao[svDB.ServidorGravacaoID] = true
}

func (s *Service) registroHandler(msg *nats.Msg) {
	reg := registro.Registro{}

	err := json.Unmarshal(msg.Data, &reg)
	if err != nil {
		s.log.Errorf("unmarshalling mensage: %s", err)
		return
	}

	_, err = s.registroCore.Create(context.Background(), reg)
	if err != nil {
		s.log.Errorf("creating registro in db: %s", err)
		s.msgr.Publish("management.processo.parar", reg.ProcessoID) // ! stop the process...
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.matchlist[reg.Placa]
	if ok {
		alerta(reg.RegistroID)
	}
}

func (s *Service) erroHandler(msg *nats.Msg) {
	s.log.Errorw("processing: %s", string(msg.Data))
}

func alerta(regID string) {
	fmt.Println("encontrado - alertaID:", regID) // TODO COLOCAR ALGUMA COISA
}
