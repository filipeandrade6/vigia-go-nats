package veiculo

import (
	"unsafe"

	"github.com/filipeandrade6/vigia-go/internal/core/veiculo/db"
	"github.com/golang/protobuf/ptypes/wrappers"
)

// TODO colcoar campos agregados e data de criacao e edicao

type Veiculo struct {
	VeiculoID string
	Placa     string
	Tipo      string
	Cor       string
	Marca     string
	Info      string
}

type NewVeiculo struct {
	Placa string `validate:"required"`
	Tipo  string `validate:"required"`
	Cor   string `validate:"required"`
	Marca string `validate:"required"`
	Info  string `validate:"required"`
}

type UpdateVeiculo struct {
	VeiculoID string                `validate:"required"`
	Placa     *wrappers.StringValue `validate:"omitempty"`
	Tipo      *wrappers.StringValue `validate:"omitempty"`
	Cor       *wrappers.StringValue `validate:"omitempty"`
	Marca     *wrappers.StringValue `validate:"omitempty"`
	Info      *wrappers.StringValue `validate:"omitempty"`
}

// =============================================================================

func toVeiculo(dbVei db.Veiculo) Veiculo {
	v := (*Veiculo)(unsafe.Pointer(&dbVei))
	return *v
}

func toVeiculoSlice(dbVeis []db.Veiculo) []Veiculo {
	veis := make([]Veiculo, len(dbVeis))
	for i, dbVei := range dbVeis {
		veis[i] = toVeiculo(dbVei)
	}

	return veis
}
