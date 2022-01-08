package db

type ServidorGravacao struct {
	ServidorGravacaoID string `db:"servidor_gravacao_id"`
	EnderecoIP         string `db:"endereco_ip"`
	Porta              int    `db:"porta"`
	Armazenamento      string `db:"armazenamento"`
	HorasRetencao      int    `db:"horas_retencao"`
}
