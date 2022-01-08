package servidorgravacao_test

import (
	"context"
	"errors"
	"testing"

	"github.com/filipeandrade6/vigia-go/internal/core/servidorgravacao"
	"github.com/filipeandrade6/vigia-go/internal/database/tests"
	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/google/go-cmp/cmp"
)

func TestServidorGravacao(t *testing.T) {
	log, db, teardown := tests.New(t)
	t.Cleanup(teardown)

	core := servidorgravacao.NewCore(log, db)

	ctx := context.Background()

	t.Log("\tGiven the need to work with Servidores de Gravacao records.")
	{
		ns := servidorgravacao.NewServidorGravacao{
			EnderecoIP:    "15.25.35.45",
			Porta:         5001,
			Armazenamento: "/home/filipe/vigia",
			HorasRetencao: 1,
		}

		sv, err := core.Create(ctx, ns)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to create servidor de gravacao: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to create servidor de gravacao.", tests.Success)

		saved, err := core.QueryByID(ctx, sv.ServidorGravacaoID)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve servidor de gravacao by ID: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve servidor de gravacao by ID.", tests.Success)

		saved2, err := core.QueryByEnderecoIPPorta(ctx, sv.EnderecoIP, sv.Porta)
		if saved2.ServidorGravacaoID != saved.ServidorGravacaoID {
			t.Fatalf("\t%s\tShould be able to retrieve servidor de gravacao by Endereco IP and Porta: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve servidor de gravacao by Endereco IP and Porta.", tests.Success)

		if diff := cmp.Diff(sv, saved); diff != "" {
			t.Fatalf("\t%s\tShould get back the same servidor de gravacao. Diff:\n%s", tests.Failed, diff)
		}
		t.Logf("\t%s\tShould get back the same servidor de gravacao.", tests.Success)

		nsv2 := servidorgravacao.NewServidorGravacao{
			EnderecoIP:    "15.25.35.45",
			Porta:         5001,
			Armazenamento: "/home/filipe/vigia",
			HorasRetencao: 2,
		}

		_, err = core.Create(ctx, nsv2)
		if !errors.Is(err, servidorgravacao.ErrServidorAlreadyExists) {
			t.Fatalf("\t%s\tShould NOT be able to create servidor de gravacao with an already existing Endereco IP and Porta: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould NOT be able to create servidor de gravacao with an already existing Endereco IP and Porta.", tests.Success)

		upd := servidorgravacao.UpdateServidorGravacao{
			ServidorGravacaoID: sv.ServidorGravacaoID,
			EnderecoIP:         &wrappers.StringValue{Value: "60.50.30.20"},
			Porta:              &wrappers.Int32Value{Value: 2727},
			Armazenamento:      &wrappers.StringValue{Value: "/home/filipe/vigia2"},
			HorasRetencao:      &wrappers.Int32Value{Value: 2},
		}

		if err = core.Update(ctx, upd); err != nil {
			t.Fatalf("\t%s\tShould be able to update servidor de gravacao: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to update servidor de gravacao.", tests.Success)

		svs, err := core.Query(ctx, "", 1, 3)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve updated servidor de gravacao: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve updated servidor de gravacao.", tests.Success)

		want := sv
		want.EnderecoIP = upd.EnderecoIP.GetValue()
		want.Porta = int(upd.Porta.GetValue())
		want.Armazenamento = upd.Armazenamento.GetValue()
		want.HorasRetencao = int(upd.HorasRetencao.GetValue())

		var idx int
		for i, s := range svs {
			if s.ServidorGravacaoID == sv.ServidorGravacaoID {
				idx = i
			}
		}
		if diff := cmp.Diff(want, svs[idx]); diff != "" {
			t.Fatalf("\t%s\tShould get back the same servidor de gravacao. Diff:\n%s", tests.Failed, diff)
		}
		t.Logf("\t%s\tShould get back the same servidor de gravacao.", tests.Success)

		upd = servidorgravacao.UpdateServidorGravacao{
			ServidorGravacaoID: sv.ServidorGravacaoID,
			Porta:              &wrappers.Int32Value{Value: 4343},
		}

		if err = core.Update(ctx, upd); err != nil {
			t.Fatalf("\t%s\tShould be able to update just some fields of servidor de gravacao: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to update just some fields of servidor de gravacao.", tests.Success)

		saved, err = core.QueryByID(ctx, sv.ServidorGravacaoID)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve servidor de gravacao by ID: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve servidor de gravacao by ID.", tests.Success)

		if saved.Porta != int(upd.Porta.GetValue()) {
			t.Fatalf("\t%s\tShould be able to see updated Porta field: got %q want %q.", tests.Failed, saved.Porta, int(upd.Porta.GetValue()))
		}
		t.Logf("\t%s\tShould be able to see updated Porta field.", tests.Success)

		if err = core.Delete(ctx, sv.ServidorGravacaoID); err != nil {
			t.Fatalf("\t%s\tShould be able to delete servidor de gravacao: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to delete servidor de gravacao.", tests.Success)

		_, err = core.QueryByID(ctx, sv.ServidorGravacaoID)
		if !errors.Is(err, servidorgravacao.ErrNotFound) {
			t.Fatalf("\t%s\tShould NOT be able to retrieve servidor de gravacao: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould NOT be able to retrieve servidor de gravacao.", tests.Success)
	}

	t.Log("\tGiven the need to page through Servidores de Gravacao records.")
	{
		sv1, err := core.Query(ctx, "", 1, 1)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve servidores de gravacao for page 1: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve servidores de gravacao for page 1.", tests.Success)

		if len(sv1) != 1 {
			t.Fatalf("\t%s\tShould have a single servidor de gravacao: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have a single servidor de gravacao.", tests.Success)

		sv2, err := core.Query(ctx, "", 2, 1)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve servidores de gravacao for page 2: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve servidores de gravacao for page 2.", tests.Success)

		if len(sv2) != 1 {
			t.Fatalf("\t%s\tShould have a single servidor de gravacao: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have a single servidor de gravacao.", tests.Success)

		if sv1[0].ServidorGravacaoID == sv2[0].ServidorGravacaoID {
			t.Logf("\t\tServidor1: %v", sv1[0].ServidorGravacaoID)
			t.Logf("\t\tServidor2: %v", sv2[0].ServidorGravacaoID)
			t.Fatalf("\t%s\tShould have different servidores de gravacao: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have different servidores de gravacao.", tests.Success)
	}
}
