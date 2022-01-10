package registro

import (
	"context"
	"errors"
	"fmt"

	"github.com/filipeandrade6/vigia-go/internal/core/registro/db"
	"github.com/filipeandrade6/vigia-go/internal/sys/database"
	"github.com/filipeandrade6/vigia-go/internal/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	ErrNotFound  = errors.New("registro not found")
	ErrInvalidID = errors.New("ID is not in its proper from")
)

type Core struct {
	store db.Store
}

func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

func (c Core) Create(ctx context.Context, r Registro) (string, error) {
	dbReg := db.Registro{
		RegistroID:    r.RegistroID,
		ProcessoID:    r.ProcessoID,
		Placa:         r.Placa,
		TipoVeiculo:   r.TipoVeiculo,
		CorVeiculo:    r.CorVeiculo,
		MarcaVeiculo:  r.MarcaVeiculo,
		Armazenamento: r.Armazenamento,
		Confianca:     r.Confianca,
		CriadoEm:      r.CriadoEm,
	}

	if err := c.store.Create(ctx, dbReg); err != nil {
		return "", fmt.Errorf("create: %w", err)
	}

	return dbReg.RegistroID, nil
}

func (c Core) Query(ctx context.Context, query string, pageNumber int, rowsPerPage int) ([]Registro, error) {
	dbRegs, err := c.store.Query(ctx, query, pageNumber, rowsPerPage)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return toRegistroSlice(dbRegs), nil
}

func (c Core) QueryByID(ctx context.Context, registroID string) (Registro, error) {
	if err := validate.CheckID(registroID); err != nil {
		return Registro{}, ErrInvalidID
	}

	dbReg, err := c.store.QueryByID(ctx, registroID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Registro{}, ErrNotFound
		}
		return Registro{}, fmt.Errorf("query: %w", err)
	}

	return toRegistro(dbReg), nil
}
