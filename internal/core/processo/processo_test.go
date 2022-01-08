package processo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/filipeandrade6/vigia-go/internal/core/processo"
	"github.com/filipeandrade6/vigia-go/internal/database/tests"

	"github.com/google/go-cmp/cmp"
)

func TestProcesso(t *testing.T) {
	log, db, teardown := tests.New(t)
	t.Cleanup(teardown)

	core := processo.NewCore(log, db)

	ctx := context.Background()

	np := processo.NewProcesso{
		ServidorGravacaoID: "d03307d4-2b28-4c23-a004-3da25e5b8bb1", // seeded
		CameraID:           "d03307d4-2b28-4c23-a004-3da25e5b8ce3", // seeded
		Processador:        2,
		Adaptador:          2,
	}

	t.Log("\tGiven the need to work with Processo records.")
	{
		prc, err := core.Create(ctx, np)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to create processo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to create processo.", tests.Success)

		saved, err := core.QueryByID(ctx, prc.ProcessoID)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve processo by ID: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve processo by ID.", tests.Success)

		if diff := cmp.Diff(prc, saved); diff != "" {
			t.Fatalf("\t%s\tShould get back the same processo. Diff:\n%s", tests.Failed, diff)
		}
		t.Logf("\t%s\tShould get back the same processo.", tests.Success)

		upd := processo.UpdateProcesso{
			ServidorGravacaoID: tests.StringPointer("d03307d4-2b28-4c23-a004-3da25e524bb1"),
			CameraID:           tests.StringPointer("d03307d4-2b28-4c23-a004-3da25e5b8aa3"),
			Processador:        tests.IntPointer(5),
			Adaptador:          tests.IntPointer(5),
		}

		if err = core.Update(ctx, prc.ProcessoID, upd); err != nil {
			t.Fatalf("\t%s\tShould be able to update processo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to update processo.", tests.Success)

		prcs, err := core.Query(ctx, "", 1, 3)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve updated processo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve updated processo.", tests.Success)

		want := prc
		want.ServidorGravacaoID = *upd.ServidorGravacaoID
		want.CameraID = *upd.CameraID
		want.Processador = *upd.Processador
		want.Adaptador = *upd.Adaptador
		// want.Execucao = *upd.Execucao

		var idx int
		for i, p := range prcs {
			if p.ProcessoID == want.ProcessoID {
				idx = i
			}
		}
		if diff := cmp.Diff(want, prcs[idx]); diff != "" {
			t.Fatalf("\t%s\tShould get back the same processo. Diff\n%s", tests.Failed, diff)
		}
		t.Logf("\t%s\tShould get back the same processo.", tests.Success)

		upd = processo.UpdateProcesso{
			Adaptador: tests.IntPointer(7),
		}

		if err = core.Update(ctx, prc.ProcessoID, upd); err != nil {
			t.Fatalf("\t%s\tShould be able to update just some fields of processo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to update just some fields of processo.", tests.Success)

		saved, err = core.QueryByID(ctx, prc.ProcessoID)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve updated processo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve updated processo.", tests.Success)

		if saved.Adaptador != *upd.Adaptador {
			t.Fatalf("\t%s\tShould be able to see updated Adaptador field: got %q want %q.", tests.Failed, saved.Adaptador, *upd.Adaptador)
		}
		t.Logf("\t%s\tShould be able to see updated Adaptador field.", tests.Success)

		if err = core.Delete(ctx, prc.ProcessoID); err != nil {
			t.Fatalf("\t%s\tShould be able to delete processo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to delete processo.", tests.Success)

		_, err = core.QueryByID(ctx, prc.ProcessoID)
		if !errors.Is(err, processo.ErrNotFound) {
			t.Fatalf("\t%s\tShould NOT be able to retrieve deleted processo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould NOT be able to retrieve deleted processo.", tests.Success)

		prcs, err = core.QueryAll(ctx)
		if err != nil || len(prcs) != 2 {
			t.Fatalf("\t%s\tShould be able to retrieve 2 processos: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve 2 processos.", tests.Success)
	}

	t.Log("\tGiven the need to page through Processo records.")
	{
		prc1, err := core.Query(ctx, "", 1, 1)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve processo for page 1: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve processo for page 1.", tests.Success)

		if len(prc1) != 1 {
			t.Fatalf("\t%s\tShould have a single processo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have a single processo.", tests.Success)

		prc2, err := core.Query(ctx, "", 2, 1)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve processo for page 2: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve processo for page 2.", tests.Success)

		if len(prc2) != 1 {
			t.Fatalf("\t%s\tShould have a single processo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have a single processo.", tests.Success)

		if prc1[0].ProcessoID == prc2[0].ProcessoID {
			t.Logf("\t\tServidor1: %v", prc1[0].ProcessoID)
			t.Logf("\t\tServidor2: %v", prc2[0].ProcessoID)
			t.Fatalf("\t%s\tShould have different processo: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have different processo.", tests.Success)
	}
}
