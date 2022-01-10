package veiculo

import (
	"context"
	"errors"
	"fmt"

	"github.com/filipeandrade6/vigia-go/internal/core/veiculo/db"
	"github.com/filipeandrade6/vigia-go/internal/sys/database"
	"github.com/filipeandrade6/vigia-go/internal/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	ErrNotFound           = errors.New("veiculo not found")
	ErrInvalidID          = errors.New("ID is not in its proper from")
	ErrPlacaAlreadyExists = errors.New("placa already exists")
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
func (c Core) Create(ctx context.Context, nv NewVeiculo) (Veiculo, error) {
	if err := validate.Check(nv); err != nil {
		return Veiculo{}, fmt.Errorf("validating data: %w", err)
	}

	if _, err := c.QueryByPlaca(ctx, nv.Placa); !errors.Is(err, ErrNotFound) {
		return Veiculo{}, ErrPlacaAlreadyExists
	}

	dbVei := db.Veiculo{
		VeiculoID: validate.GenerateID(),
		Placa:     nv.Placa,
		Tipo:      nv.Tipo,
		Cor:       nv.Cor,
		Marca:     nv.Marca,
		Info:      nv.Info,
	}

	if err := c.store.Create(ctx, dbVei); err != nil {
		return Veiculo{}, fmt.Errorf("create: %w", err)
	}

	return toVeiculo(dbVei), nil
}

func (c Core) Update(ctx context.Context, up UpdateVeiculo) error {
	if err := validate.CheckID(up.VeiculoID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(up); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbVei, err := c.store.QueryByID(ctx, up.VeiculoID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating veiculo veiculoID[%s]: %w", up.VeiculoID, err)
	}

	if up.Placa != nil {
		if _, err := c.QueryByPlaca(ctx, up.Placa.GetValue()); !errors.Is(err, ErrNotFound) {
			return ErrPlacaAlreadyExists
		}
		dbVei.Placa = up.Placa.GetValue()
	}
	if up.Tipo != nil {
		dbVei.Tipo = up.Tipo.GetValue()
	}
	if up.Cor != nil {
		dbVei.Cor = up.Cor.GetValue()
	}
	if up.Marca != nil {
		dbVei.Marca = up.Marca.GetValue()
	}
	if up.Info != nil {
		dbVei.Info = up.Info.GetValue()
	}

	if err := c.store.Update(ctx, dbVei); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (c Core) Delete(ctx context.Context, veiculoID string) error {
	if err := validate.CheckID(veiculoID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, veiculoID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c Core) Query(ctx context.Context, query string, pageNumber int, rowsPerPage int) ([]Veiculo, error) {
	dbVeis, err := c.store.Query(ctx, query, pageNumber, rowsPerPage)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return toVeiculoSlice(dbVeis), nil
}

func (c Core) QueryByID(ctx context.Context, veiculoID string) (Veiculo, error) {
	if err := validate.CheckID(veiculoID); err != nil {
		return Veiculo{}, ErrInvalidID
	}

	dbVei, err := c.store.QueryByID(ctx, veiculoID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Veiculo{}, ErrNotFound
		}
		return Veiculo{}, fmt.Errorf("query: %w", err)
	}

	return toVeiculo(dbVei), nil
}

func (c Core) QueryByPlaca(ctx context.Context, placa string) (Veiculo, error) {
	dbVei, err := c.store.QueryByPlaca(ctx, placa)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Veiculo{}, ErrNotFound
		}
		return Veiculo{}, fmt.Errorf("query: %w", err)
	}

	return toVeiculo(dbVei), nil
}

func (c Core) QueryAll(ctx context.Context) ([]Veiculo, error) {
	dbVeis, err := c.store.QueryAll(ctx)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return toVeiculoSlice(dbVeis), nil
}
