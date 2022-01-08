package veiculo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/filipeandrade6/vigia-go/internal/core/veiculo"
	"github.com/filipeandrade6/vigia-go/internal/database/tests"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/go-cmp/cmp"
)

// TODO teste com email repetido

func TestVeiculo(t *testing.T) {
	log, db, teardown := tests.New(t)
	t.Cleanup(teardown)

	core := veiculo.NewCore(log, db)

	ctx := context.Background()

	t.Log("\tGiven the need to work with Veiculo records.")
	{
		nv := veiculo.NewVeiculo{
			Placa: "XYZ0000",
			Tipo:  "sedan",
			Cor:   "vermelho",
			Marca: "fiat",
			Info:  "teste de informacao",
		}

		vei, err := core.Create(ctx, nv)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to create veiculo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to create veiculo.", tests.Success)

		saved, err := core.QueryByID(ctx, vei.VeiculoID)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve veiculo by ID: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve veiculo by ID.", tests.Success)

		saved2, err := core.QueryByPlaca(ctx, vei.Placa)
		if saved2.VeiculoID != saved.VeiculoID {
			t.Fatalf("\t%s\tShould be able to retrieve veiculo by Placa: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve veiculo by Placa.", tests.Success)

		if diff := cmp.Diff(vei, saved); diff != "" {
			t.Fatalf("\t%s\tShould get back the same veiculo. Diff:\n%s", tests.Failed, diff)
		}
		t.Logf("\t%s\tShould get back the same veiculo.", tests.Success)

		nu2 := veiculo.NewVeiculo{
			Placa: "XYZ0000",
			Tipo:  "suv",
			Cor:   "preto",
			Marca: "fiat",
			Info:  "teste de informacao",
		}

		_, err = core.Create(ctx, nu2)
		if !errors.Is(err, veiculo.ErrPlacaAlreadyExists) {
			t.Fatalf("\t%s\tShould NOT be able to create veiculo with an already existing Placa: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould NOT be able to create veiculo with an already existing Placa.", tests.Success)

		upd := veiculo.UpdateVeiculo{
			VeiculoID: vei.VeiculoID,
			Placa:     &wrappers.StringValue{Value: "ABC1111"},
			Tipo:      &wrappers.StringValue{Value: "suv"},
			Cor:       &wrappers.StringValue{Value: "preto"},
			Marca:     &wrappers.StringValue{Value: "fiat"},
			Info:      &wrappers.StringValue{Value: "teste de informacao"},
		}

		if err = core.Update(ctx, upd); err != nil {
			t.Fatalf("\t%s\tShould be able to update veiculo: %s.", tests.Failed, err) // TODO dando erro aqui de email ja existe, mas não existe
		}
		t.Logf("\t%s\tShould be able to update veiculo.", tests.Success)

		veis, err := core.Query(ctx, "", 1, 3)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve Veiculos: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve Veiculos.", tests.Success)

		want := vei
		want.Placa = upd.Placa.GetValue()
		want.Tipo = upd.Tipo.GetValue()
		want.Cor = upd.Cor.GetValue()
		want.Marca = upd.Marca.GetValue()
		want.Info = upd.Info.GetValue()

		var idx int
		for i, v := range veis {
			if v.VeiculoID == vei.VeiculoID {
				idx = i
			}
		}
		if diff := cmp.Diff(want, veis[idx]); diff != "" {
			t.Fatalf("\t%s\tShould get back the same servidor de gravacao. Diff:\n%s", tests.Failed, diff)
		}
		t.Logf("\t%s\tShould get back the same servidor de gravacao.", tests.Success)

		upd = veiculo.UpdateVeiculo{
			VeiculoID: vei.VeiculoID,
			Marca:     &wrappers.StringValue{Value: "bmw"},
		}

		if err = core.Update(ctx, upd); err != nil {
			t.Fatalf("\t%s\tShould be able to update just some fields of veiculo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to update just some fields of veiculo.", tests.Success)

		saved, err = core.QueryByID(ctx, vei.VeiculoID)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve veiculo by Email: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve veiculo by Email.", tests.Success)

		if saved.Marca != upd.Marca.GetValue() {
			t.Fatalf("\t%s\tShould be able to see updated Email field: got %q want %q.", tests.Failed, saved.Marca, upd.Marca.GetValue())
		}
		t.Logf("\t%s\tShould be able to see updated Email field.", tests.Success)

		if err := core.Delete(ctx, vei.VeiculoID); err != nil {
			t.Fatalf("\t%s\tShould be able to delete veiculo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to delete veiculo.", tests.Success)

		_, err = core.QueryByID(ctx, vei.VeiculoID)
		if !errors.Is(err, veiculo.ErrNotFound) {
			t.Fatalf("\t%s\tShould NOT be able to retrieve veiculo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould NOT be able to retrieve veiculo.", tests.Success)

		veis, err = core.QueryAll(ctx)
		if err != nil || len(veis) != 2 {
			t.Fatalf("\t%s\tShould be able to retrieve 2 veiculos: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve 2 veiculos.", tests.Success)
	}

	t.Log("\tGiven the need to page through Veiculo records.")
	{
		vei1, err := core.Query(ctx, "", 1, 1)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve users for page 1: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve users for page 1.", tests.Success)

		if len(vei1) != 1 {
			t.Fatalf("\t%s\tShould have a single user : %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have a single user.", tests.Success)

		vei2, err := core.Query(ctx, "", 2, 1)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve users for page 2: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve users for page 2.", tests.Success)

		if len(vei2) != 1 {
			t.Fatalf("\t%s\tShould have a single user: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have a single user.", tests.Success)

		if vei1[0].VeiculoID == vei2[0].VeiculoID {
			t.Logf("\t\tUser1: %v", vei1[0].VeiculoID)
			t.Logf("\t\tUser2: %v", vei2[0].VeiculoID)
			t.Fatalf("\t%s\tShould have different users: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have different users.", tests.Success)
	}
}
