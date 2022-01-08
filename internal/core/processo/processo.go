package processo

import (
	"context"
	"errors"
	"fmt"

	"github.com/filipeandrade6/vigia-go/internal/core/processo/db"
	"github.com/filipeandrade6/vigia-go/internal/sys/database"
	"github.com/filipeandrade6/vigia-go/internal/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// TODO atualizar isso aqui quando chegar a hora

var (
	ErrNotFound  = errors.New("processo not found")
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

// TODO entender o validator...
func (c Core) Create(ctx context.Context, np NewProcesso) (Processo, error) {
	if err := validate.Check(np); err != nil {
		return Processo{}, fmt.Errorf("validating data: %w", err)
	}

	dbPrc := db.Processo{
		ProcessoID:         validate.GenerateID(),
		ServidorGravacaoID: np.ServidorGravacaoID,
		CameraID:           np.CameraID,
		Processador:        np.Processador,
		Adaptador:          np.Adaptador,
	}

	if err := c.store.Create(ctx, dbPrc); err != nil {
		return Processo{}, fmt.Errorf("create: %w", err)
	}

	return toProcesso(dbPrc), nil
}

func (c Core) Update(ctx context.Context, processoID string, up UpdateProcesso) error {
	if err := validate.CheckID(processoID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(up); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbPrc, err := c.store.QueryByID(ctx, processoID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating processo processoID[%s]: %w", processoID, err)
	}

	if up.ServidorGravacaoID != nil {
		dbPrc.ServidorGravacaoID = *up.ServidorGravacaoID
	}
	if up.CameraID != nil {
		dbPrc.CameraID = *up.CameraID
	}
	if up.Processador != nil {
		dbPrc.Processador = *up.Processador
	}
	if up.Adaptador != nil {
		dbPrc.Adaptador = *up.Adaptador
	}

	if err := c.store.Update(ctx, dbPrc); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (c Core) Delete(ctx context.Context, processoID string) error {
	if err := validate.CheckID(processoID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, processoID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c Core) Query(ctx context.Context, query string, pageNumber int, rowsPerPage int) ([]Processo, error) {
	dbPrcs, err := c.store.Query(ctx, query, pageNumber, rowsPerPage)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return toProcessoSlice(dbPrcs), nil
}

func (c Core) QueryByID(ctx context.Context, processoID string) (Processo, error) {
	if err := validate.CheckID(processoID); err != nil {
		return Processo{}, ErrInvalidID
	}

	dbPrc, err := c.store.QueryByID(ctx, processoID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Processo{}, ErrNotFound
		}
		return Processo{}, fmt.Errorf("query: %w", err)
	}

	return toProcesso(dbPrc), nil
}

func (c Core) QueryAll(ctx context.Context) ([]Processo, error) {
	dbPrcs, err := c.store.QueryAll(ctx)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return toProcessoSlice(dbPrcs), nil
}
