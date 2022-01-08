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

func (s Store) Create(ctx context.Context, vei Veiculo) error {
	const q = `
	INSERT INTO veiculos
		(veiculo_id, placa, tipo, cor, marca, info)
	VALUES
		(:veiculo_id, :placa, :tipo, :cor, :marca, :info)`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, vei); err != nil {
		return fmt.Errorf("inserting veiculo: %w", err)
	}

	return nil
}

func (s Store) Update(ctx context.Context, vei Veiculo) error {
	const q = `
	UPDATE
		veiculos
	SET
		"placa" = :placa,
		"tipo" = :tipo,
		"cor" = :cor,
		"marca" = :marca,
		"info" = :info
	WHERE
		veiculo_id = :veiculo_id`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, vei); err != nil {
		return fmt.Errorf("updating veiculo veiculoID[%s]: %w", vei.VeiculoID, err)
	}

	return nil
}

func (s Store) Delete(ctx context.Context, veiculoID string) error {
	data := struct {
		VeiculoID string `db:"veiculo_id"`
	}{
		VeiculoID: veiculoID,
	}

	const q = `
	DELETE FROM
		veiculos
	WHERE
		veiculo_id = :veiculo_id`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, data); err != nil {
		return fmt.Errorf("deleting veiculo veiculoID[%s]: %w", veiculoID, err)
	}

	return nil
}

func (s Store) Query(ctx context.Context, query string, pageNumber int, rowsPerPage int) ([]Veiculo, error) {
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
		veiculos
	WHERE
		CONCAT(veiculo_id, placa, tipo, cor, marca, info)
	ILIKE
		:query
	ORDER BY
		placa
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var veiculos []Veiculo
	if err := database.NamedQuerySlice(ctx, s.log, s.sqlxDB, q, data, &veiculos); err != nil {
		return nil, fmt.Errorf("selecting veiculos [%s]: %w", query, err)
	}

	return veiculos, nil
}

func (s Store) QueryByID(ctx context.Context, veiculoID string) (Veiculo, error) {
	data := struct {
		VeiculoID string `db:"veiculo_id"`
	}{
		VeiculoID: veiculoID,
	}

	const q = `
	SELECT
		*
	FROM
		veiculos
	WHERE
		veiculo_id = :veiculo_id`

	var vei Veiculo
	if err := database.NamedQueryStruct(ctx, s.log, s.sqlxDB, q, data, &vei); err != nil {
		return Veiculo{}, fmt.Errorf("selecting veiculo veiculoID[%q]: %w", veiculoID, err)
	}

	return vei, nil
}

func (s Store) QueryByPlaca(ctx context.Context, placa string) (Veiculo, error) {
	data := struct {
		Placa string `db:"placa"`
	}{
		Placa: placa,
	}

	const q = `
	SELECT
		*
	FROM
		veiculos
	WHERE
		placa = :placa`

	var vei Veiculo
	if err := database.NamedQueryStruct(ctx, s.log, s.sqlxDB, q, data, &vei); err != nil {
		return Veiculo{}, fmt.Errorf("selecting veiculo placa[%q]: %w", placa, err)
	}

	return vei, nil
}

func (s Store) QueryAll(ctx context.Context) ([]Veiculo, error) {
	const q = `
	SELECT
		*
	FROM
		veiculos`

	var veis []Veiculo
	if err := database.QuerySlice(ctx, s.log, s.sqlxDB, q, &veis); err != nil {
		return nil, fmt.Errorf("selecting veiculos: %w", err)
	}

	return veis, nil
}
