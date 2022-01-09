package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/filipeandrade6/vigia-go/internal/core/registro"
	"github.com/filipeandrade6/vigia-go/internal/core/servidorgravacao"
	"github.com/nats-io/nats.go"
)

// Se descober retornar nula é pq tem erro
func (s *Service) descobertaHandler(msg *nats.Msg) {
	// ? ver se vai funcionar isso aqui
	var err error
	defer func(err error) {
		if err != nil {
			if err = msg.Respond(nil); err != nil {
				s.log.Errorf("responding mensage: %s", err) // ! disparar um alerta
			}
		}
	}(err)

	var sv servidorgravacao.NewServidorGravacao

	err = json.Unmarshal(msg.Data, &sv)
	if err != nil {
		s.log.Errorf("unmarshalling mensage: %s", err) // ! disparar um alerta
		return
	}

	svDB, err := s.servidorGravacaoCore.Create(context.Background(), sv)
	if err != nil {
		s.log.Errorf("creating servidor gravacao in database: %s", err) // ! disparar um alerta
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
		s.log.Errorf("unmarshalling mensage: %s", err) // ! disparar um alerta
		return
	}

	_, err = s.registroCore.Create(context.Background(), reg)
	if err != nil {
		s.log.Errorf("creating registro in db: %s", err)            // ! disparar um alerta
		s.msgr.Publish("management.processo.parar", reg.ProcessoID) // stop the process...
		return
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.matchlist[reg.Placa]
	if ok {
		alerta(reg.RegistroID)
	}
}

func (s *Service) infoHandler(msg *nats.Msg) {
	fmt.Println("info handler acionado") // TODO COLOCAR ALGUMA COSIA
}

func (s *Service) erroHandler(msg *nats.Msg) {
	fmt.Println("error handler acionado") // TODO COLOCAR ALGUMA COSIA
}

func alerta(regID string) {
	fmt.Println("encontrado - alertaID:", regID) // TODO COLOCAR ALGUMA COISA
}
