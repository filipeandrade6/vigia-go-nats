package veiculo

import (
	"unsafe"

	pb "github.com/filipeandrade6/vigia-go/internal/api/v1"
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

// =============================================================================

func (v Veiculo) ToProto() *pb.Veiculo {
	return &pb.Veiculo{
		VeiculoId: v.VeiculoID,
		Placa:     v.Placa,
		Tipo:      v.Tipo,
		Cor:       v.Cor,
		Marca:     v.Marca,
		Info:      v.Info,
	}
}

func FromProto(v *pb.Veiculo) Veiculo {
	return Veiculo{
		VeiculoID: v.GetVeiculoId(),
		Placa:     v.GetPlaca(),
		Tipo:      v.GetTipo(),
		Cor:       v.GetCor(),
		Marca:     v.GetMarca(),
		Info:      v.GetInfo(),
	}
}

type Veiculos []Veiculo

func (v Veiculos) ToProto() []*pb.Veiculo {
	var veis []*pb.Veiculo

	for _, vei := range v {
		veis = append(veis, vei.ToProto())
	}

	return veis
}

func VeiculosFromProto(v []*pb.Veiculo) Veiculos {
	var veis Veiculos

	for _, vei := range v {
		veis = append(veis, FromProto(vei))
	}

	return veis
}
