package db

import "github.com/lib/pq"

type Usuario struct {
	UsuarioID string         `db:"usuario_id"`
	Email     string         `db:"email"`
	Funcao    pq.StringArray `db:"funcao"`
	Senha     []byte         `db:"senha"`
}
