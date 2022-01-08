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

func (s Store) Create(ctx context.Context, prc Processo) error {
	const q = `
	INSERT INTO processos
		(processo_id, servidor_gravacao_id, camera_id, processador, adaptador)
	VALUES
		(:processo_id, :servidor_gravacao_id, :camera_id, :processador, :adaptador)`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, prc); err != nil {
		return fmt.Errorf("inserting processo: %w", err)
	}

	return nil
}

func (s Store) Update(ctx context.Context, prc Processo) error {
	const q = `
	UPDATE
		processos
	SET
		"servidor_gravacao_id" = :servidor_gravacao_id,
		"camera_id" = :camera_id,
		"processador" = :processador,
		"adaptador" = :adaptador
	WHERE
		processo_id = :processo_id`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, prc); err != nil {
		return fmt.Errorf("updating processo processoID[%s]: %w", prc.ProcessoID, err)
	}

	return nil
}

func (s Store) Delete(ctx context.Context, processoID string) error {
	data := struct {
		ProcessoID string `db:"processo_id"`
	}{
		ProcessoID: processoID,
	}

	const q = `
	DELETE FROM
		processos
	WHERE
		processo_id = :processo_id`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, data); err != nil {
		return fmt.Errorf("deleting processo processoID[%s]: %w", processoID, err)
	}

	return nil
}

func (s Store) Query(ctx context.Context, query string, pageNumber int, rowsPerPage int) ([]Processo, error) {
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
		processos
	WHERE
		CONCAT(processo_id, servidor_gravacao_id, camera_id, processador, adaptador)
	ILIKE
		:query
	ORDER BY
		servidor_gravacao_id
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var prcs []Processo
	if err := database.NamedQuerySlice(ctx, s.log, s.sqlxDB, q, data, &prcs); err != nil {
		return nil, fmt.Errorf("selecting processos [%s]: %w", query, err)
	}

	return prcs, nil
}

func (s Store) QueryByID(ctx context.Context, processoID string) (Processo, error) {
	data := struct {
		ProcessoID string `db:"processo_id"`
	}{
		ProcessoID: processoID,
	}

	const q = `
	SELECT
		*
	FROM
		processos
	WHERE
		processo_id = :processo_id`

	var prc Processo
	if err := database.NamedQueryStruct(ctx, s.log, s.sqlxDB, q, data, &prc); err != nil {
		return Processo{}, fmt.Errorf("selecting processo processoID[%q]: %w", processoID, err)
	}

	return prc, nil
}

func (s Store) QueryAll(ctx context.Context) ([]Processo, error) {
	const q = `
	SELECT
		*
	FROM
		processos`

	var prcs []Processo
	if err := database.QuerySlice(ctx, s.log, s.sqlxDB, q, &prcs); err != nil {
		return nil, fmt.Errorf("selecting processos: %w", err)
	}

	return prcs, nil
}
