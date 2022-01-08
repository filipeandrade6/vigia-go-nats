package db

import "time"

type Registro struct {
	RegistroID    string    `db:"registro_id"`
	ProcessoID    string    `db:"processo_id"`
	Placa         string    `db:"placa"`
	TipoVeiculo   string    `db:"tipo_veiculo"`
	CorVeiculo    string    `db:"cor_veiculo"`
	MarcaVeiculo  string    `db:"marca_veiculo"`
	Armazenamento string    `db:"armazenamento"`
	Confianca     float32   `db:"confianca"`
	CriadoEm      time.Time `db:"criado_em"`
}
