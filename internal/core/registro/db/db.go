package db

import (
	"context"
	"fmt"

	"github.com/filipeandrade6/vigia-go/internal/sys/database"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Store struct {
	log    *zap.SugaredLogger
	sqlxDB *sqlx.DB
}

func NewStore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Store {
	return Store{
		log:    log,
		sqlxDB: sqlxDB,
	}
}

func (s Store) Create(ctx context.Context, reg Registro) error {
	const q = `
	INSERT INTO registros
		(registro_id, processo_id, placa, tipo_veiculo, cor_veiculo, marca_veiculo, armazenamento, confianca, criado_em)
	VALUES
		(:registro_id, :processo_id, :placa, :tipo_veiculo, :cor_veiculo, :marca_veiculo, :armazenamento, :confianca, :criado_em)`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, reg); err != nil {
		return fmt.Errorf("inserting registro: %w", err)
	}

	return nil
}

func (s Store) Query(ctx context.Context, query string, pageNumber int, rowsPerPage int) ([]Registro, error) {
	data := struct {
		Query       string `db:"query"`
		Offset      int    `db:"offset"`
		RowsPerPage int    `db:"rows_per_page"`
	}{
		Query:       fmt.Sprintf("%%%s%%", query),
		Offset:      (pageNumber - 1) * rowsPerPage,
		RowsPerPage: rowsPerPage,
	}

	const q = `
	SELECT
		*
	FROM
		registros
	WHERE
		CONCAT(registro_id, processo_id, placa, tipo_veiculo, cor_veiculo, marca_veiculo, armazenamento, confianca, criado_em)
	ILIKE
		:query
	ORDER BY
		criado_em
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var cams []Registro
	if err := database.NamedQuerySlice(ctx, s.log, s.sqlxDB, q, data, &cams); err != nil {
		return nil, fmt.Errorf("selecting registros [%s]: %w", query, err)
	}

	return cams, nil
}

func (s Store) QueryByID(ctx context.Context, registroID string) (Registro, error) {
	data := struct {
		RegistroID string `db:"registro_id"`
	}{
		RegistroID: registroID,
	}

	const q = `
	SELECT
		*
	FROM
		registros
	WHERE
		registro_id = :registro_id`

	var reg Registro
	if err := database.NamedQueryStruct(ctx, s.log, s.sqlxDB, q, data, &reg); err != nil {
		return Registro{}, fmt.Errorf("selecting registro registroID[%q]: %w", registroID, err)
	}

	return reg, nil
}
