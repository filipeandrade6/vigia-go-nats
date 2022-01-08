package operrors

import "errors"

// type OpErr int

// const (
// 	InitNetSDK = iota
// 	Unauthorized
// 	Unreachable
// 	BadRequest
// 	NetSDKLogin
// 	SaveImage
// 	Analyzer
// 	NotDahua
// )

// func (o OpErr) String() string {
// 	return [...]string{
// 		"Unauthorized",
// 		"Unreachable",
// 		"BadRequest",
// 		"NetSDKLogin",
// 		"SaveImage",
// 		"Analyzer",
// 		"NotDahua",
// 	}[o]
// }

var (
	ErrInitNetSDK   = errors.New("init netsdk failed")
	ErrUnauthorized = errors.New("unauthorized")
	ErrUnreachable  = errors.New("unreachable")
	ErrBadRequest   = errors.New("bad request")
	ErrNetSDKLogin  = errors.New("netsdk login failed")
	ErrSaveImage    = errors.New("could not save image")
	ErrAnalyzer     = errors.New("analyzer error")
	ErrNotDahua     = errors.New("not dahua")
)

type OpError struct {
	ServidorGravacaoID string
	ProcessoID         string
	RegistroID         string
	// Err        OpErr
	Err             error
	StoppedProcesso bool
}

func (o *OpError) Unwrap() error { return o.Err } // TODO: verificar isso aqui

func (o *OpError) Error() string {
	// if o == nil {
	// 	return "<nil>"
	// }

	s := "srv[" + o.ServidorGravacaoID + "]"
	if o.ProcessoID != "" {
		s += " prc[" + o.ProcessoID + "]"
	}
	if o.RegistroID != "" {
		s += " reg[" + o.RegistroID + "]"
	}
	// s += ": " + o.Err.String()
	s += ": " + o.Err.Error()

	return s
}
