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

func (s Store) Create(ctx context.Context, sv ServidorGravacao) error {
	const q = `
	INSERT INTO servidores_gravacao
		(servidor_gravacao_id, endereco_ip, porta, armazenamento, horas_retencao)
	VALUES
		(:servidor_gravacao_id, :endereco_ip, :porta, :armazenamento, :horas_retencao)`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, sv); err != nil {
		return fmt.Errorf("inserting servidor de gravacao: %w", err)
	}

	return nil
}

func (s Store) Update(ctx context.Context, sv ServidorGravacao) error {
	const q = `
	UPDATE
		servidores_gravacao
	SET
		"endereco_ip" = :endereco_ip,
		"porta" = :porta,
		"armazenamento" = :armazenamento,
		"horas_retencao" = :horas_retencao
	WHERE
		servidor_gravacao_id = :servidor_gravacao_id`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, sv); err != nil {
		return fmt.Errorf("updating servidor de gravacao svID[%s]: %w", sv.ServidorGravacaoID, err)
	}

	return nil
}

func (s Store) Delete(ctx context.Context, svID string) error {
	data := struct {
		SvID string `db:"servidor_gravacao_id"`
	}{
		SvID: svID,
	}

	const q = `
	DELETE FROM
		servidores_gravacao
	WHERE
		servidor_gravacao_id = :servidor_gravacao_id`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, data); err != nil {
		return fmt.Errorf("deleting servidor de gravacao svID[%s]: %w", svID, err)
	}

	return nil
}

func (s Store) Query(ctx context.Context, query string, pageNumber int, rowsPerPage int) ([]ServidorGravacao, error) {
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
		servidores_gravacao
	WHERE
		CONCAT(servidor_gravacao_id, endereco_ip, porta, armazenamento, horas_retencao)
	ILIKE
		:query
	ORDER BY
		endereco_ip
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var svs []ServidorGravacao
	if err := database.NamedQuerySlice(ctx, s.log, s.sqlxDB, q, data, &svs); err != nil {
		return nil, fmt.Errorf("selecting servidor de gravacao [%s]: %w", query, err)
	}

	return svs, nil
}

func (s Store) QueryByID(ctx context.Context, svID string) (ServidorGravacao, error) {
	data := struct {
		SvID string `db:"servidor_gravacao_id"`
	}{
		SvID: svID,
	}

	const q = `
	SELECT
		*
	FROM
		servidores_gravacao
	WHERE
		servidor_gravacao_id = :servidor_gravacao_id`

	var sv ServidorGravacao
	if err := database.NamedQueryStruct(ctx, s.log, s.sqlxDB, q, data, &sv); err != nil {
		return ServidorGravacao{}, fmt.Errorf("selecting servidor de gravacao svID[%q]: %w", svID, err)
	}

	return sv, nil
}

func (s Store) QueryByEnderecoIPPorta(ctx context.Context, endereco_ip string, porta int) (ServidorGravacao, error) {
	data := struct {
		EnderecoIP string `db:"endereco_ip"`
		Porta      int    `db:"porta"`
	}{
		EnderecoIP: endereco_ip,
		Porta:      porta,
	}

	const q = `
	SELECT
		*
	FROM
		servidores_gravacao
	WHERE
		endereco_ip = :endereco_ip AND porta = :porta`

	var sv ServidorGravacao
	if err := database.NamedQueryStruct(ctx, s.log, s.sqlxDB, q, data, &sv); err != nil {
		return ServidorGravacao{}, fmt.Errorf("selecting servidor de gravacao enderecoIP:Porta[%q:%q]: %w", endereco_ip, porta, err)
	}

	return sv, nil
}
