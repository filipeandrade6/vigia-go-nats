package camera

import (
	"unsafe"

	"github.com/filipeandrade6/vigia-go/internal/core/camera/db"
	"github.com/golang/protobuf/ptypes/wrappers"
)

// TODO colcoar campos agregados e data de criacao e edicao

type Camera struct {
	CameraID   string `json:"camera_id"`
	Descricao  string `json:"descricao"`
	EnderecoIP string `json:"endereco_ip"`
	Porta      int    `json:"porta"`
	Canal      int    `json:"canal"`
	Usuario    string `json:"usuario"`
	Senha      string `json:"senha"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
}

type NewCamera struct {
	Descricao  string `validate:"required"`
	EnderecoIP string `validate:"required,ip"`
	Porta      int    `validate:"required,gte=1,lte=65536"`
	Canal      int    `validate:"required,gte=0,lte=10"`
	Usuario    string `validate:"required"`
	Senha      string `validate:"required"`
	Latitude   string `validate:"required,latitude"`
	Longitude  string `validate:"required,longitude"`
}

type UpdateCamera struct {
	CameraID   string                `validate:"required"`
	Descricao  *wrappers.StringValue `validate:"omitempty"`
	EnderecoIP *wrappers.StringValue `validate:"omitempty,ip"`
	Porta      *wrappers.Int32Value  `validate:"omitempty,gte=1,lte=65536"`
	Canal      *wrappers.Int32Value  `validate:"omitempty,gte=0,lte=10"`
	Usuario    *wrappers.StringValue `validate:"omitempty"`
	Senha      *wrappers.StringValue `validate:"omitempty"`
	Latitude   *wrappers.StringValue `validate:"omitempty,latitude"`
	Longitude  *wrappers.StringValue `validate:"omitempty,longitude"`
}

// =============================================================================

func toCamera(dbCam db.Camera) Camera {
	c := (*Camera)(unsafe.Pointer(&dbCam))
	return *c
}

func toCameraSlice(dbCams []db.Camera) []Camera {
	cams := make([]Camera, len(dbCams))
	for i, dbCam := range dbCams {
		cams[i] = toCamera(dbCam)
	}
	return cams
}
