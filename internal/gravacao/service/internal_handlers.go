package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/filipeandrade6/vigia-go/internal/core/registro"
	"github.com/filipeandrade6/vigia-go/internal/core/servidorgravacao"
	"github.com/nats-io/nats.go"
)

func (s *Service) descobertaHandler(msg *nats.Msg) {
	var sv servidorgravacao.NewServidorGravacao
	if err := json.Unmarshal(msg.Data, &sv); err != nil {
		fmt.Println(err) // ! tratar
	}

	svDB, err := s.servidorGravacaoCore.Create(context.Background(), sv)
	if err != nil {
		fmt.Println(err) // ! tratar
	}

	if err = msg.Respond([]byte(svDB.ServidorGravacaoID)); err != nil {
		fmt.Println(err) // ! tratar
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.servidoresGravacao[svDB.ServidorGravacaoID] = true
}

func (s *Service) registroHandler(msg *nats.Msg) {
	reg := registro.Registro{}
	err := json.Unmarshal(msg.Data, &reg)
	if err != nil {
		fmt.Println(err) // ! tratar
	}

	_, err = s.registroCore.Create(context.Background(), reg)
	if err != nil {
		fmt.Println(err) // ! tratar
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
