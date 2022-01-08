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

func (s Store) Create(ctx context.Context, cam Camera) error {
	const q = `
	INSERT INTO cameras
		(camera_id, descricao, endereco_ip, porta, canal, usuario, senha, latitude, longitude)
	VALUES
		(:camera_id, :descricao, :endereco_ip, :porta, :canal, :usuario, :senha, :latitude, :longitude)`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, cam); err != nil {
		return fmt.Errorf("inserting camera: %w", err)
	}

	return nil
}

func (s Store) Update(ctx context.Context, cam Camera) error {
	const q = `
	UPDATE
		cameras
	SET
		"descricao" = :descricao,
		"endereco_ip" = :endereco_ip,
		"porta" = :porta,
		"canal" = :canal,
		"usuario" = :usuario,
		"senha" = :senha,
		"latitude" = :latitude,
		"longitude" = :longitude
	WHERE
		camera_id = :camera_id`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, cam); err != nil {
		return fmt.Errorf("updating camera cameraID[%s]: %w", cam.CameraID, err)
	}

	return nil
}

func (s Store) Delete(ctx context.Context, cameraID string) error {
	data := struct {
		CameraID string `db:"camera_id"`
	}{
		CameraID: cameraID,
	}

	const q = `
	DELETE FROM
		cameras
	WHERE
		camera_id = :camera_id`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, data); err != nil {
		return fmt.Errorf("deleting camera cameraID[%s]: %w", cameraID, err)
	}

	return nil
}

func (s Store) Query(ctx context.Context, query string, pageNumber int, rowsPerPage int) ([]Camera, error) {
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
		cameras
	WHERE
		CONCAT(camera_id, descricao, endereco_ip, porta, canal, usuario, senha, latitude, longitude)
	ILIKE
		:query
	ORDER BY
		descricao
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var cams []Camera
	if err := database.NamedQuerySlice(ctx, s.log, s.sqlxDB, q, data, &cams); err != nil {
		return nil, fmt.Errorf("selecting cameras [%s]: %w", query, err)
	}

	return cams, nil
}

func (s Store) QueryByID(ctx context.Context, cameraID string) (Camera, error) {
	data := struct {
		CameraID string `db:"camera_id"`
	}{
		CameraID: cameraID,
	}

	const q = `
	SELECT
		*
	FROM
		cameras
	WHERE
		camera_id = :camera_id`

	var cam Camera
	if err := database.NamedQueryStruct(ctx, s.log, s.sqlxDB, q, data, &cam); err != nil {
		return Camera{}, fmt.Errorf("selecting camera cameraID[%q]: %w", cameraID, err)
	}

	return cam, nil
}

func (s Store) QueryByEnderecoIP(ctx context.Context, endereco_ip string) (Camera, error) {
	data := struct {
		EnderecoIP string `db:"endereco_ip"`
	}{
		EnderecoIP: endereco_ip,
	}

	const q = `
	SELECT
		*
	FROM
		cameras
	WHERE
		endereco_ip = :endereco_ip`

	var cam Camera
	if err := database.NamedQueryStruct(ctx, s.log, s.sqlxDB, q, data, &cam); err != nil {
		return Camera{}, fmt.Errorf("selecting camera enderecoIP[%q]: %w", endereco_ip, err)
	}

	return cam, nil
}
