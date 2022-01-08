package processo

import (
	"unsafe"

	pb "github.com/filipeandrade6/vigia-go/internal/api/v1"
	"github.com/filipeandrade6/vigia-go/internal/core/processo/db"
)

type Processo struct {
	ProcessoID         string
	ServidorGravacaoID string
	CameraID           string
	Processador        int
	Adaptador          int
}

type NewProcesso struct {
	ServidorGravacaoID string `validate:"required"`
	CameraID           string `validate:"required"`
	Processador        int    `validate:"required"`
	Adaptador          int    `validate:"required"`
}

type UpdateProcesso struct {
	ServidorGravacaoID *string `validate:"omitempty"`
	CameraID           *string `validate:"omitempty"`
	Processador        *int    `validate:"omitempty"`
	Adaptador          *int    `validate:"omitempty"`
}

// =============================================================================

func toProcesso(dbPrc db.Processo) Processo {
	p := (*Processo)(unsafe.Pointer(&dbPrc))
	return *p
}

func toProcessoSlice(dbPrcs []db.Processo) []Processo {
	prcs := make([]Processo, len(dbPrcs))
	for i, dbPrc := range dbPrcs {
		prcs[i] = toProcesso(dbPrc)
	}
	return prcs
}

// =============================================================================

func (p Processo) ToProto() *pb.Processo {
	return &pb.Processo{
		ProcessoId:         p.ProcessoID,
		ServidorGravacaoId: p.ServidorGravacaoID,
		CameraId:           p.CameraID,
		Processador:        int32(p.Processador),
		Adaptador:          int32(p.Adaptador),
	}
}

func FromProto(p *pb.Processo) Processo {
	return Processo{
		ProcessoID:         p.GetProcessoId(),
		ServidorGravacaoID: p.GetServidorGravacaoId(),
		CameraID:           p.GetCameraId(),
		Processador:        int(p.GetProcessador()),
		Adaptador:          int(p.GetAdaptador()),
	}
}

type Processos []Processo

func (p Processos) ToProto() []*pb.Processo {
	var prcs []*pb.Processo

	for _, prc := range p {
		prcs = append(prcs, prc.ToProto())
	}

	return prcs
}
