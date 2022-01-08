package db

type Processo struct {
	ProcessoID         string `db:"processo_id"`
	ServidorGravacaoID string `db:"servidor_gravacao_id"`
	CameraID           string `db:"camera_id"`
	Processador        int    `db:"processador"`
	Adaptador          int    `db:"adaptador"`
}
