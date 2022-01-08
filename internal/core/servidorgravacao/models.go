package servidorgravacao

import (
	"unsafe"

	pb "github.com/filipeandrade6/vigia-go/internal/api/v1"
	"github.com/filipeandrade6/vigia-go/internal/core/servidorgravacao/db"
	"github.com/golang/protobuf/ptypes/wrappers"
)

type ServidorGravacao struct {
	ServidorGravacaoID string
	EnderecoIP         string
	Porta              int
	Armazenamento      string
	HorasRetencao      int
}

type NewServidorGravacao struct {
	EnderecoIP    string `validate:"required,ip"`
	Porta         int    `validate:"required,gte=1,lte=65536"`
	Armazenamento string `validate:"required"`
	HorasRetencao int    `validate:"required"`
}

type UpdateServidorGravacao struct {
	ServidorGravacaoID string                `validate:"required"`
	EnderecoIP         *wrappers.StringValue `validate:"omitempty,ip"`
	Porta              *wrappers.Int32Value  `validate:"omitempty,gte=1,lte=65536"`
	Armazenamento      *wrappers.StringValue `validate:"omitempty"`
	HorasRetencao      *wrappers.Int32Value  `validate:"omitempty"`
}

// =============================================================================

func toServidorGravacao(dbSV db.ServidorGravacao) ServidorGravacao {
	s := (*ServidorGravacao)(unsafe.Pointer(&dbSV))
	return *s
}

func toServidorGravacaoSlice(dbSVs []db.ServidorGravacao) []ServidorGravacao {
	svs := make([]ServidorGravacao, len(dbSVs))
	for i, dbSV := range dbSVs {
		svs[i] = toServidorGravacao(dbSV)
	}
	return svs
}

// =============================================================================

func (s ServidorGravacao) ToProto() *pb.ServidorGravacao {
	return &pb.ServidorGravacao{
		ServidorGravacaoId: s.ServidorGravacaoID,
		EnderecoIp:         s.EnderecoIP,
		Porta:              int32(s.Porta),
		Armazenamento:      s.Armazenamento,
		HorasRetencao:      int32(s.HorasRetencao),
	}
}

func FromProto(s *pb.ServidorGravacao) ServidorGravacao {
	return ServidorGravacao{
		ServidorGravacaoID: s.GetServidorGravacaoId(),
		EnderecoIP:         s.GetEnderecoIp(),
		Porta:              int(s.GetPorta()),
		Armazenamento:      s.GetArmazenamento(),
		HorasRetencao:      int(s.GetHorasRetencao()),
	}
}

type ServidoresGravacao []ServidorGravacao

func (s ServidoresGravacao) ToProto() []*pb.ServidorGravacao {
	var svs []*pb.ServidorGravacao

	for _, sv := range s {
		svs = append(svs, sv.ToProto())
	}

	return svs
}

func ServidoresGravacaoFromProto(s []*pb.ServidorGravacao) ServidoresGravacao {
	var svs ServidoresGravacao

	for _, sv := range s {
		svs = append(svs, FromProto(sv))
	}

	return svs
}
