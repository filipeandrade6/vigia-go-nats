package registro

import (
	"time"
	"unsafe"

	pb "github.com/filipeandrade6/vigia-go/internal/api/v1"
	"github.com/filipeandrade6/vigia-go/internal/core/registro/db"
)

// TODO colcoar campos agregados e data de criacao e edicao

type Registro struct {
	RegistroID    string    `json:"registro_id"`
	ProcessoID    string    `json:"processo_id"`
	Placa         string    `json:"placa"`
	TipoVeiculo   string    `json:"tipo_veiculo"`
	CorVeiculo    string    `json:"cor_veiculo"`
	MarcaVeiculo  string    `json:"marca_veiculo"`
	Armazenamento string    `json:"armazenamento"`
	Confianca     float32   `json:"confianca"`
	CriadoEm      time.Time `json:"criado_em"`
}

// =============================================================================

func toRegistro(dbReg db.Registro) Registro {
	r := (*Registro)(unsafe.Pointer(&dbReg))
	return *r
}

func toRegistroSlice(dbRegs []db.Registro) []Registro {
	regs := make([]Registro, len(dbRegs))
	for i, dbReg := range dbRegs {
		regs[i] = toRegistro(dbReg)
	}
	return regs
}

// =============================================================================

func (r Registro) ToProto() *pb.Registro {
	return &pb.Registro{
		RegistroId: r.RegistroID,
	}
}

func FromProto(r *pb.Registro) Registro {
	return Registro{
		RegistroID: r.GetRegistroId(),
	}
}

type Registros []Registro

func (r Registros) ToProto() []*pb.Registro {
	var regs []*pb.Registro

	for _, reg := range r {
		regs = append(regs, reg.ToProto())
	}

	return regs
}

func RegistrosFromProto(r []*pb.Registro) Registros { // TODO ver se esta sendo utilizado
	var regs Registros

	for _, reg := range r {
		regs = append(regs, FromProto(reg))
	}

	return regs
}
