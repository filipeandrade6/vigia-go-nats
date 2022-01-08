package db

type Veiculo struct {
	VeiculoID string `db:"veiculo_id"`
	Placa     string `db:"placa"`
	Tipo      string `db:"tipo"`
	Cor       string `db:"cor"`
	Marca     string `db:"marca"`
	Info      string `db:"info"`
}
