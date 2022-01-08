package db

type Camera struct {
	CameraID   string `db:"camera_id"`
	Descricao  string `db:"descricao"`
	EnderecoIP string `db:"endereco_ip"`
	Porta      int    `db:"porta"`
	Canal      int    `db:"canal"`
	Usuario    string `db:"usuario"`
	Senha      string `db:"senha"`
	Latitude   string `db:"latitude"`
	Longitude  string `db:"longitude"`
}
