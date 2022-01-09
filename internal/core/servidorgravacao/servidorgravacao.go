package servidorgravacao

import (
	"context"
	"errors"
	"fmt"

	"github.com/filipeandrade6/vigia-go/internal/core/servidorgravacao/db"
	"github.com/filipeandrade6/vigia-go/internal/sys/database"
	"github.com/filipeandrade6/vigia-go/internal/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	ErrNotFound              = errors.New("servidor de gravacao not found")
	ErrInvalidID             = errors.New("ID is not in its proper from")
	ErrServidorAlreadyExists = errors.New("servidor de gravacao com endereco_ip:porta already exists")
)

type Core struct {
	store db.Store
}

func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

func (c Core) Create(ctx context.Context, nsv NewServidorGravacao) (ServidorGravacao, error) {
	if err := validate.Check(nsv); err != nil {
		return ServidorGravacao{}, fmt.Errorf("validating data: %w", err)
	}

	if _, err := c.QueryByEnderecoIPPorta(ctx, nsv.EnderecoIP, nsv.Porta); !errors.Is(err, ErrNotFound) {
		return ServidorGravacao{}, ErrServidorAlreadyExists
	}

	dbSV := db.ServidorGravacao{
		ServidorGravacaoID: validate.GenerateID(),
		EnderecoIP:         nsv.EnderecoIP,
		Porta:              nsv.Porta,
	}

	if err := c.store.Create(ctx, dbSV); err != nil {
		return ServidorGravacao{}, fmt.Errorf("create: %w", err)
	}

	return toServidorGravacao(dbSV), nil
}

func (c Core) Update(ctx context.Context, up UpdateServidorGravacao) error {
	if err := validate.CheckID(up.ServidorGravacaoID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(up); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbSV, err := c.store.QueryByID(ctx, up.ServidorGravacaoID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating servidor de gravacao svID[%s]: %w", up.ServidorGravacaoID, err)
	}

	if up.EnderecoIP != nil {
		dbSV.EnderecoIP = up.EnderecoIP.GetValue()
	}
	if up.Porta != nil {
		dbSV.Porta = int(up.Porta.GetValue())
	}
	if up.Armazenamento != nil {
		dbSV.Armazenamento = up.Armazenamento.GetValue()
	}
	if up.HorasRetencao != nil {
		dbSV.HorasRetencao = int(up.HorasRetencao.GetValue())
	}

	if _, err := c.QueryByEnderecoIPPorta(ctx, dbSV.EnderecoIP, dbSV.Porta); !errors.Is(err, ErrNotFound) {
		return ErrServidorAlreadyExists
	}

	if err := c.store.Update(ctx, dbSV); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (c Core) Delete(ctx context.Context, svID string) error {
	if err := validate.CheckID(svID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, svID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// TODO arrumar
// func (c Core) Query(ctx context.Context, query string, pageNumber int, rowsPerPage int) (ServidoresGravacao, error) {
// 	dbSVs, err := c.store.Query(ctx, query, pageNumber, rowsPerPage)
// 	if err != nil {
// 		if errors.Is(err, database.ErrDBNotFound) {
// 			return nil, ErrNotFound
// 		}
// 		return nil, fmt.Errorf("query: %w", err)
// 	}

// 	return toServidorGravacaoSlice(dbSVs), nil
// }

func (c Core) QueryByID(ctx context.Context, svID string) (ServidorGravacao, error) {
	if err := validate.CheckID(svID); err != nil {
		return ServidorGravacao{}, ErrInvalidID
	}

	dbSV, err := c.store.QueryByID(ctx, svID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ServidorGravacao{}, ErrNotFound
		}
		return ServidorGravacao{}, fmt.Errorf("query: %w", err)
	}

	return toServidorGravacao(dbSV), nil
}

func (c Core) QueryByEnderecoIPPorta(ctx context.Context, endereco_ip string, porta int) (ServidorGravacao, error) {
	dbSv, err := c.store.QueryByEnderecoIPPorta(ctx, endereco_ip, porta)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ServidorGravacao{}, ErrNotFound
		}
		return ServidorGravacao{}, fmt.Errorf("query: %w", err)
	}

	return toServidorGravacao(dbSv), nil
}
