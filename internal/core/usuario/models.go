package usuario

import (
	"unsafe"

	"github.com/filipeandrade6/vigia-go/internal/core/usuario/db"

	wrappers "github.com/golang/protobuf/ptypes/wrappers"
)

type Usuario struct {
	UsuarioID string
	Email     string
	Funcao    []string
	Senha     string
}

type NewUsuario struct {
	Email  string   `validate:"required,email"`
	Funcao []string `validate:"required"`
	Senha  string   `validate:"required"`
}

type UpdateUsuario struct {
	UsuarioID string                `validate:"required"`
	Email     *wrappers.StringValue `validate:"omitempty,email"`
	Funcao    []string              `validate:"omitempty"`
	Senha     *wrappers.StringValue `validate:"omitempty"`
}

// =============================================================================

func toUsuario(dbUsr db.Usuario) Usuario {
	pu := (*Usuario)(unsafe.Pointer(&dbUsr))
	return *pu
}

func toUsuarioSlice(dbUsrs []db.Usuario) []Usuario {
	usuarios := make([]Usuario, len(dbUsrs))
	for i, dbUsr := range dbUsrs {
		usuarios[i] = toUsuario(dbUsr)
	}
	return usuarios
}
