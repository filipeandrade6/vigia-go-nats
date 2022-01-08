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

func (s Store) Create(ctx context.Context, usuario Usuario) error {
	const q = `
	INSERT INTO usuarios
		(usuario_id, email, senha, funcao)
	VALUES
		(:usuario_id, :email, :senha, :funcao)`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, usuario); err != nil {
		return fmt.Errorf("inserting usuario: %w", err)
	}

	return nil
}

func (s Store) Update(ctx context.Context, usuario Usuario) error {
	const q = `
	UPDATE
		usuarios
	SET
		"email" = :email,
		"senha" = :senha,
		"funcao" = :funcao
	WHERE
		usuario_id = :usuario_id`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, usuario); err != nil {
		return fmt.Errorf("updating usuario usuarioID[%s]: %w", usuario.UsuarioID, err)
	}

	return nil
}

func (s Store) Delete(ctx context.Context, usuarioID string) error {
	data := struct {
		UsuarioID string `db:"usuario_id"`
	}{
		UsuarioID: usuarioID,
	}

	const q = `
	DELETE FROM
		usuarios
	WHERE
		usuario_id = :usuario_id`

	if err := database.NamedExecContext(ctx, s.log, s.sqlxDB, q, data); err != nil {
		return fmt.Errorf("deleting usuario usuarioID[%s]: %w", usuarioID, err)
	}

	return nil
}

func (s Store) Query(ctx context.Context, query string, pageNumber int, rowsPerPage int) ([]Usuario, error) {
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
		usuarios
	WHERE
		CONCAT(usuario_id, email, senha, funcao)
	ILIKE
		:query
	ORDER BY
		email
	OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY`

	var usuarios []Usuario
	if err := database.NamedQuerySlice(ctx, s.log, s.sqlxDB, q, data, &usuarios); err != nil {
		return nil, fmt.Errorf("selecting usuarios [%s]: %w", query, err)
	}

	return usuarios, nil
}

func (s Store) QueryByID(ctx context.Context, usuarioID string) (Usuario, error) {
	data := struct {
		UsuarioID string `db:"usuario_id"`
	}{
		UsuarioID: usuarioID,
	}

	const q = `
	SELECT
		*
	FROM
		usuarios
	WHERE
		usuario_id = :usuario_id`

	var usuario Usuario
	if err := database.NamedQueryStruct(ctx, s.log, s.sqlxDB, q, data, &usuario); err != nil {
		return Usuario{}, fmt.Errorf("selecting usuario usuarioID[%q]: %w", usuarioID, err)
	}

	return usuario, nil
}

func (s Store) QueryByEmail(ctx context.Context, email string) (Usuario, error) {
	data := struct {
		Email string `db:"email"`
	}{
		Email: email,
	}

	const q = `
	SELECT
		*
	FROM
		usuarios
	WHERE
		email = :email`

	var usuario Usuario
	if err := database.NamedQueryStruct(ctx, s.log, s.sqlxDB, q, data, &usuario); err != nil {
		return Usuario{}, fmt.Errorf("selecting usuario email[%q]: %w", email, err)
	}

	return usuario, nil
}

// func (s Store) Authenticate(ctx context.Context, email, senha string) (auth.Claims, error) {
// 	data := struct {
// 		Email string `db:"email"`
// 	}{
// 		Email: email,
// 	}

// 	fmt.Println("chegou aqui")

// 	const q = `
// 	SELECT
// 		*
// 	FROM
// 		usuarios
// 	WHERE
// 		email = :email`

// 	var usuario Usuario
// 	if err := database.NamedQueryStruct(ctx, s.log, s.sqlxDB, q, data, &usuario); err != nil {
// 		if errors.As(err, &database.ErrNotFound) {
// 			return auth.Claims{}, database.ErrNotFound
// 		}
// 		return auth.Claims{}, fmt.Errorf("selecting usuario[%q]: %w", email, err)
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Senha), []byte(senha)); err != nil {
// 		return auth.Claims{}, database.ErrAuthenticationFailure
// 	}

// 	claims := auth.Claims{
// 		StandardClaims: jwt.StandardClaims{
// 			Issuer:    "service project",
// 			Subject:   usuario.UsuarioID,
// 			ExpiresAt: time.Now().Add(time.Hour).Unix(),
// 			IssuedAt:  time.Now().UTC().Unix(),
// 		},
// 		Roles: usuario.Funcao,
// 	}

// 	return claims, nil
// }
