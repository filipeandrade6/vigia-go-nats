package usuario

import (
	"context"
	"errors"
	"fmt"

	"github.com/filipeandrade6/vigia-go/internal/core/usuario/db"
	"github.com/filipeandrade6/vigia-go/internal/sys/database"
	"github.com/filipeandrade6/vigia-go/internal/sys/validate"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("usuario not found")
	ErrInvalidID             = errors.New("ID is not in its proper form")
	ErrAuthenticationFailure = errors.New("authentication failed") // TODO ver se no service tem metodo login aqui
	ErrEmailAlreadyExists    = errors.New("email already exists")
)

type Core struct {
	store db.Store
}

func NewCore(log *zap.SugaredLogger, sqlxDB *sqlx.DB) Core {
	return Core{
		store: db.NewStore(log, sqlxDB),
	}
}

func (c Core) Create(ctx context.Context, nu NewUsuario) (Usuario, error) {
	if err := validate.Check(nu); err != nil {
		return Usuario{}, fmt.Errorf("validating data: %w", err)
	}

	if _, err := c.QueryByEmail(ctx, nu.Email); !errors.Is(err, ErrNotFound) {
		return Usuario{}, ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Senha), bcrypt.DefaultCost)
	if err != nil {
		return Usuario{}, fmt.Errorf("generating password hash: %w", err)
	}

	dbUsr := db.Usuario{
		UsuarioID: validate.GenerateID(),
		Email:     nu.Email,
		Funcao:    nu.Funcao,
		Senha:     hash,
	}

	if err := c.store.Create(ctx, dbUsr); err != nil {
		return Usuario{}, fmt.Errorf("create: %w", err)
	}

	return toUsuario(dbUsr), nil
}

func (c Core) Update(ctx context.Context, up UpdateUsuario) error {
	if err := validate.CheckID(up.UsuarioID); err != nil {
		return ErrInvalidID
	}

	if err := validate.Check(up); err != nil {
		return fmt.Errorf("validating data: %w", err)
	}

	dbUsr, err := c.store.QueryByID(ctx, up.UsuarioID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("updating user userID[%s]: %w", up.UsuarioID, err)
	}

	if up.Email != nil {
		if _, err := c.QueryByEmail(ctx, up.Email.GetValue()); !errors.Is(err, ErrNotFound) {
			return ErrEmailAlreadyExists
		}
		dbUsr.Email = up.Email.GetValue()
	}
	if up.Funcao != nil {
		dbUsr.Funcao = up.Funcao
	}
	if up.Senha != nil {
		pw, err := bcrypt.GenerateFromPassword([]byte(up.Senha.GetValue()), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("generating password hash: %w", err)
		}
		dbUsr.Senha = pw
	}

	if err := c.store.Update(ctx, dbUsr); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (c Core) Delete(ctx context.Context, usuarioID string) error {
	if err := validate.CheckID(usuarioID); err != nil {
		return ErrInvalidID
	}

	if err := c.store.Delete(ctx, usuarioID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func (c Core) Query(ctx context.Context, query string, pageNumber int, rowsPerPage int) ([]Usuario, error) {
	dbUsrs, err := c.store.Query(ctx, query, pageNumber, rowsPerPage)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	return toUsuarioSlice(dbUsrs), nil
}

func (c Core) QueryByID(ctx context.Context, usuarioID string) (Usuario, error) {
	if err := validate.CheckID(usuarioID); err != nil {
		return Usuario{}, database.ErrInvalidID
	}

	dbUsr, err := c.store.QueryByID(ctx, usuarioID)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Usuario{}, ErrNotFound
		}
		return Usuario{}, fmt.Errorf("query: %w", err)
	}

	return toUsuario(dbUsr), nil
}

func (c Core) QueryByEmail(ctx context.Context, email string) (Usuario, error) {
	dbUsr, err := c.store.QueryByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, database.ErrDBNotFound) {
			return Usuario{}, ErrNotFound
		}
		return Usuario{}, fmt.Errorf("query: %w", err)
	}

	return toUsuario(dbUsr), nil
}

// func (c Core) Authenticate(ctx context.Context, email, senha string) (auth.Claims, error) {
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
// 	if err := database.NamedQueryStruct(ctx, s.log, s.db, q, data, &usuario); err != nil {
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
