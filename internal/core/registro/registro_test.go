package registro_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/filipeandrade6/vigia-go/internal/core/registro"
	"github.com/filipeandrade6/vigia-go/internal/database/tests"
	"github.com/filipeandrade6/vigia-go/internal/sys/validate"

	"github.com/google/go-cmp/cmp"
)

func TestRegistro(t *testing.T) {
	log, db, teardown := tests.New(t)
	t.Cleanup(teardown)

	core := registro.NewCore(log, db)

	ctx := context.Background()

	t.Log("\tGiven the need to work with Registro records.")
	{
		now := time.Now()

		nr := registro.Registro{
			RegistroID:   validate.GenerateID(),
			ProcessoID:   "d03307d4-2b28-4c23-a004-3da32e5b8a61",
			Placa:        "XXX1111",
			TipoVeiculo:  "sedan",
			CorVeiculo:   "prata",
			MarcaVeiculo: "honda",
			Confianca:    0.50,
			CriadoEm:     now,
		}

		reg, err := core.Create(ctx, nr)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to create registro: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to create registro.", tests.Success)

		saved, err := core.QueryByID(ctx, reg)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve registro by ID: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve registro by ID.", tests.Success)

		nr.RegistroID = reg
		nr.Armazenamento = fmt.Sprintf("%d_%s", now.Unix(), reg)

		want := nr
		want.Armazenamento = saved.Armazenamento
		want.CriadoEm = saved.CriadoEm

		if diff := cmp.Diff(want, saved); diff != "" {
			t.Fatalf("\t%s\tShould get back the same registro. Diff:\n%s", tests.Failed, diff)
		}
		t.Logf("\t%s\tShould get back the same registro.", tests.Success)
	}

	t.Log("\tGiven the need to page through Registro records.")
	{
		reg1, err := core.Query(ctx, "", 1, 1)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve registro for page 1: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve registro for page 1.", tests.Success)

		if len(reg1) != 1 {
			t.Fatalf("\t%s\tShould have a single registro: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have a single registro.", tests.Success)

		reg2, err := core.Query(ctx, "", 2, 1)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve registro for page 2: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve registro for page 2.", tests.Success)

		if len(reg2) != 1 {
			t.Fatalf("\t%s\tShould have a single registro: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have a single registro.", tests.Success)

		if reg1[0].RegistroID == reg2[0].RegistroID {
			t.Logf("\t\tRegistro1: %v", reg1[0].RegistroID)
			t.Logf("\t\tRegistro2: %v", reg2[0].RegistroID)
			t.Fatalf("\t%s\tShould have different registro: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have different registro.", tests.Success)
	}
}
