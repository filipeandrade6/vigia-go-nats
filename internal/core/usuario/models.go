package usuario

import (
	"unsafe"

	pb "github.com/filipeandrade6/vigia-go/internal/api/v1"
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

// =============================================================================

func (u Usuario) ToProto() *pb.Usuario {
	return &pb.Usuario{
		UsuarioId: u.UsuarioID,
		Email:     u.Email,
		Funcao:    u.Funcao,
		Senha:     u.Senha,
	}
}

func FromProto(u *pb.Usuario) Usuario {
	return Usuario{
		UsuarioID: u.GetUsuarioId(),
		Email:     u.GetEmail(),
		Funcao:    u.GetFuncao(),
		Senha:     u.GetSenha(),
	}
}

type Usuarios []Usuario

func (u Usuarios) ToProto() []*pb.Usuario {
	var usuarios []*pb.Usuario

	for _, usuario := range u {
		usuarios = append(usuarios, usuario.ToProto())
	}

	return usuarios
}

func UsuariosFromProto(u []*pb.Usuario) Usuarios { // TODO ver se esta sendo utilizado
	var usrs Usuarios

	for _, usr := range u {
		usrs = append(usrs, FromProto(usr))
	}

	return usrs
}
