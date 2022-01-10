package camera

import (
	"context"
	"errors"
	"fmt"

	"github.com/filipeandrade6/vigia-go/internal/core/camera/db"
	"github.com/filipeandrade6/vigia-go/internal/sys/database"
	"github.com/filipeandrade6/vigia-go/internal/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	ErrNotFound            = errors.New("camera not found")
	ErrInvalidID           = errors.New("ID is not in its proper from")
	ErrCameraAlreadyExists = errors.New("camera already exists")
)

type Core struct {
	store db.Store
}

func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

// TODO entender o validator...
func (c Core) Create(ctx context.Context, nc NewCamera) (Camera, error) {
	if err := validate.Check(nc); err != nil {
		return Camera{}, fmt.Errorf("validating data: %w", err)
	}

	if _, err := c.QueryByEnderecoIP(ctx, nc.EnderecoIP); !errors.Is(err, ErrNotFound) {
		return Camera{}, ErrCameraAlreadyExists
	}

	dbCam := db.Camera{
		CameraID:   validate.GenerateID(),
		Descricao:  nc.Descricao,
		EnderecoIP: nc.EnderecoIP,
		Porta:      nc.Porta,
		Canal:      nc.Canal,
		Usuario:    nc.Usuario,
		Senha:      nc.Senha,
		Latitude:   nc.Latitude,
		Longitude:  nc.Longitude,
	}

	if err := c.store.Create(ctx, dbCam); err != nil {
		return Camera{}, fmt.Errorf("create: %w", err)
	}

	return toCamera(dbCam), nil
}

func (c Core) Update(ctx context.Context, up UpdateCamera) error {
	if err := validate.CheckID(up.CameraID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(up); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbCam, err := c.store.QueryByID(ctx, up.CameraID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating camera cameraID[%s]: %w", up.CameraID, err)
	}

	if up.Descricao != nil {
		dbCam.Descricao = up.Descricao.GetValue()
	}
	if up.EnderecoIP != nil {
		if _, err := c.QueryByEnderecoIP(ctx, up.EnderecoIP.GetValue()); !errors.Is(err, ErrNotFound) {
			return ErrCameraAlreadyExists
		}
		dbCam.EnderecoIP = up.EnderecoIP.GetValue()
	}
	if up.Porta != nil {
		dbCam.Porta = int(up.Porta.GetValue())
	}
	if up.Canal != nil {
		dbCam.Canal = int(up.Canal.GetValue())
	}
	if up.Usuario != nil {
		dbCam.Usuario = up.Usuario.GetValue()
	}
	if up.Senha != nil {
		dbCam.Senha = up.Senha.GetValue()
	}
	if up.Latitude != nil {
		dbCam.Latitude = up.Latitude.GetValue()
	}
	if up.Longitude != nil {
		dbCam.Longitude = up.Longitude.GetValue()
	}

	if err := c.store.Update(ctx, dbCam); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (c Core) Delete(ctx context.Context, cameraID string) error {
	if err := validate.CheckID(cameraID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, cameraID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c Core) Query(ctx context.Context, query string, pageNumber int, rowsPerPage int) ([]Camera, error) {
	dbCams, err := c.store.Query(ctx, query, pageNumber, rowsPerPage)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return toCameraSlice(dbCams), nil
}

func (c Core) QueryByID(ctx context.Context, cameraID string) (Camera, error) {
	if err := validate.CheckID(cameraID); err != nil {
		return Camera{}, ErrInvalidID
	}

	dbCam, err := c.store.QueryByID(ctx, cameraID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Camera{}, ErrNotFound
		}
		return Camera{}, fmt.Errorf("query: %w", err)
	}

	return toCamera(dbCam), nil
}

func (c Core) QueryByEnderecoIP(ctx context.Context, enderecoIP string) (Camera, error) {
	dbCam, err := c.store.QueryByEnderecoIP(ctx, enderecoIP)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Camera{}, ErrNotFound
		}
		return Camera{}, fmt.Errorf("query: %w", err)
	}

	return toCamera(dbCam), nil
}
