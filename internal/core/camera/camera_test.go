package camera_test

import (
	"context"
	"errors"
	"testing"

	"github.com/filipeandrade6/vigia-go/internal/core/camera"
	"github.com/filipeandrade6/vigia-go/internal/database/tests"
	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/google/go-cmp/cmp"
)

func TestCamera(t *testing.T) {
	log, db, teardown := tests.New(t)
	t.Cleanup(teardown)

	core := camera.NewCore(log, db)

	ctx := context.Background()

	t.Log("\tGiven the need to work with Camera records.")
	{
		nc := camera.NewCamera{
			Descricao:  "camera testes 1",
			EnderecoIP: "1.2.3.4",
			Porta:      1234,
			Canal:      1,
			Usuario:    "admin",
			Senha:      "admin",
			Latitude:   "-12.4567",
			Longitude:  "-12.4567",
		}

		cam, err := core.Create(ctx, nc)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to create camera: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to create camera.", tests.Success)

		saved, err := core.QueryByID(ctx, cam.CameraID)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve camera by ID: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve camera by ID.", tests.Success)

		saved2, err := core.QueryByEnderecoIP(ctx, cam.EnderecoIP)
		if saved2.CameraID != saved.CameraID {
			t.Fatalf("\t%s\tShould be able to retrieve camera by Endereco IP: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve camera by Endereco IP.", tests.Success)

		if diff := cmp.Diff(cam, saved); diff != "" {
			t.Fatalf("\t%s\tShould get back the same camera. Diff:\n%s", tests.Failed, diff)
		}
		t.Logf("\t%s\tShould get back the same camera.", tests.Success)

		nc2 := camera.NewCamera{
			Descricao:  "camera testes 2",
			EnderecoIP: "1.2.3.4",
			Porta:      1235,
			Canal:      2,
			Usuario:    "manager",
			Senha:      "manager",
			Latitude:   "-13.4567",
			Longitude:  "-13.4567",
		}

		_, err = core.Create(ctx, nc2)
		if !errors.Is(err, camera.ErrCameraAlreadyExists) {
			t.Fatalf("\t%s\tShould NOT be able to create camera with an already existing Endereco IP: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould NOT be able to create camera with an already existing Endereco IP.", tests.Success)

		upd := camera.UpdateCamera{
			CameraID:   cam.CameraID,
			Descricao:  &wrappers.StringValue{Value: "camera atualizada"},
			EnderecoIP: &wrappers.StringValue{Value: "123.123.210.210"},
			Porta:      &wrappers.Int32Value{Value: 2020},
			Canal:      &wrappers.Int32Value{Value: 7},
			Usuario:    &wrappers.StringValue{Value: "user"},
			Senha:      &wrappers.StringValue{Value: "user"},
			Latitude:   &wrappers.StringValue{Value: "-10.1010"},
			Longitude:  &wrappers.StringValue{Value: "-10.1010"},
		}

		if err = core.Update(ctx, upd); err != nil {
			t.Fatalf("\t%s\tShould be able to update camera: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to update camera.", tests.Success)

		cams, err := core.Query(ctx, "", 1, 3)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve cameras: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve cameras.", tests.Success)

		want := cam
		want.Descricao = upd.Descricao.GetValue()
		want.EnderecoIP = upd.EnderecoIP.GetValue()
		want.Porta = int(upd.Porta.GetValue())
		want.Canal = int(upd.Canal.GetValue())
		want.Usuario = upd.Usuario.GetValue()
		want.Senha = upd.Senha.GetValue()
		want.Latitude = upd.Latitude.GetValue()
		want.Longitude = upd.Longitude.GetValue()

		var idx int
		for i, c := range cams {
			if c.CameraID == want.CameraID {
				idx = i
			}
		}
		if diff := cmp.Diff(want, cams[idx]); diff != "" {
			t.Fatalf("\t%s\tShould get back the same camera. Diff\n%s", tests.Failed, diff)
		}
		t.Logf("\t%s\tShould get back the same camera.", tests.Success)

		upd = camera.UpdateCamera{
			CameraID: cam.CameraID,
			Porta:    &wrappers.Int32Value{Value: 5321},
		}

		if err = core.Update(ctx, upd); err != nil {
			t.Fatalf("\t%s\tShould be able to update just some fields of camera: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to update just some fields of camera.", tests.Success)

		saved, err = core.QueryByID(ctx, cam.CameraID)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve updated camera: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve updated camera.", tests.Success)

		if saved.Porta != int(upd.Porta.GetValue()) {
			t.Fatalf("\t%s\tShould be able to see updated Porta field: got %q want %q.", tests.Failed, saved.Porta, int(upd.Porta.GetValue()))
		}
		t.Logf("\t%s\tShould be able to see updated Porta field.", tests.Success)

		if err = core.Delete(ctx, cam.CameraID); err != nil {
			t.Fatalf("\t%s\tShould be able to delete camera: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to delete camera.", tests.Success)

		_, err = core.QueryByID(ctx, cam.CameraID)
		if !errors.Is(err, camera.ErrNotFound) {
			t.Fatalf("\t%s\tShould NOT be able to retrieve deleted camera: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould NOT be able to retrieve deleted camera.", tests.Success)
	}

	t.Log("\tGiven the need to page through Camera records.")
	{
		cam1, err := core.Query(ctx, "", 1, 1)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve camera for page 1: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve camera for page 1.", tests.Success)

		if len(cam1) != 1 {
			t.Fatalf("\t%s\tShould have a single camera: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have a single camera.", tests.Success)

		cam2, err := core.Query(ctx, "", 2, 1)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to retrieve camera for page 2: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to retrieve camera for page 2.", tests.Success)

		if len(cam2) != 1 {
			t.Fatalf("\t%s\tShould have a single camera: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have a single camera.", tests.Success)

		if cam1[0].CameraID == cam2[0].CameraID {
			t.Logf("\t\tServidor1: %v", cam1[0].CameraID)
			t.Logf("\t\tServidor2: %v", cam2[0].CameraID)
			t.Fatalf("\t%s\tShould have different camera: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould have different camera.", tests.Success)
	}
}
