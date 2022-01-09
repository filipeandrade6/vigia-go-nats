package servidorgravacao

import (
	"unsafe"

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
	EnderecoIP    string `json:"endereco_ip" validate:"required,ip"`
	Porta         int    `json:"porta" validate:"required,gte=1,lte=65536"`
	Armazenamento string `json:"armazenamento" validate:"required"`
	HorasRetencao int    `json:"horas_retencao" validate:"required"`
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

// func toServidorGravacaoSlice(dbSVs []db.ServidorGravacao) []ServidorGravacao {
// 	svs := make([]ServidorGravacao, len(dbSVs))
// 	for i, dbSV := range dbSVs {
// 		svs[i] = toServidorGravacao(dbSV)
// 	}
// 	return svs
// }
