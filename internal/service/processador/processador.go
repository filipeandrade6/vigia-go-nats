package processador

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/filipeandrade6/vigia-go/internal/core/registro"
	"github.com/filipeandrade6/vigia-go/internal/sys/operrors"
)

// TODO utilizar defer nos mux pois melhorar a leitura > performance.

type Camera interface {
	// New(processoID, enderecoIP string, porta, canal int, usuario, senha string)
	Start(armazenamento string, regChan chan registro.Registro, errChan chan *operrors.OpError)
	Stop()
	GetID() string
}

type Processador struct {
	armazenamento string
	horasRetencao int

	registroCore registro.Core
	errChan      chan operrors.OpError // TODO usar ponteiros?
	matchChan    chan string

	mu        *sync.RWMutex
	processos map[string]Camera // TODO era ponteiros
	retry     map[string]Camera // TODO era ponteiros
	matchlist map[string]bool

	interErrChan chan *operrors.OpError
	regChan      chan registro.Registro
}

func New(
	armazenamento string,
	horasRetencao int,
	registroCore registro.Core,
	errChan chan operrors.OpError,
	matchChan chan string,
) *Processador {
	return &Processador{
		armazenamento: armazenamento,
		horasRetencao: horasRetencao,

		registroCore: registroCore,
		errChan:      errChan,
		matchChan:    matchChan,

		mu:        &sync.RWMutex{},
		processos: make(map[string]Camera),
		retry:     make(map[string]Camera),
		matchlist: make(map[string]bool),

		interErrChan: make(chan *operrors.OpError),
		regChan:      make(chan registro.Registro),
	}
}

// =================================================================================
// Processador

func (p *Processador) Start() {
	tickerHK := time.NewTicker(time.Hour)
	tickerRetry := time.NewTicker(30 * time.Second)

	for {
		select {
		// TODO ver qual o tipo de erro que da quando a camera estiver conectada e ficar offline...
		case err := <-p.interErrChan:
			switch {
			case errors.As(err.Err, &operrors.ErrUnreachable):
				p.retry[err.ProcessoID] = p.processos[err.ProcessoID]
				delete(p.processos, err.ProcessoID)

			default: // TODO ver isso abaixo
				err.StoppedProcesso = true // ! isso aqui esta em uso?
				delete(p.processos, err.ProcessoID)
			}

			p.errChan <- *err // TODO usar ponteiros?

		case <-tickerHK.C:
			go p.begintHousekeeper()

		case <-tickerRetry.C:
			for processoID, processo := range p.retry {
				p.mu.Lock()
				p.processos[processoID] = processo
				p.mu.Unlock()

				processo.Start(p.armazenamento, p.regChan, p.interErrChan)
			}
		}
	}
}

func (p *Processador) Stop() error {
	var prc []string
	p.mu.RLock()
	for k := range p.processos {
		prc = append(prc, k)
	}
	for k := range p.retry {
		prc = append(prc, k)
	}
	p.mu.RUnlock()

	// if nonStoppedPrc := p.StopProcessos(prc); nonStoppedPrc != nil {
	// 	return fmt.Errorf("could not stop processos: %s", nonStoppedPrc)
	// }

	return nil
}

// =================================================================================
// Armazenamento

func (p *Processador) UpdateArmazenamento(armazenamento string, horasRetencao int) error {
	fmt.Print("ola")
	return nil
}

// 	prcsBkp := make(map[string]Camera)
// 	p.mu.RLock()
// 	for k, v := range p.processos {
// 		prcsBkp[k] = v
// 	}

// 	var prcs []string
// 	for k := range prcsBkp {
// 		prcs = append(prcs, k)
// 	}
// 	p.mu.RUnlock()

// 	if nonStoppedPrc := p.StopProcessos(prcs); nonStoppedPrc != nil {
// 		return fmt.Errorf("could not stop processos: %s", nonStoppedPrc)
// 	}

// 	p.mu.Lock()
// 	p.armazenamento = armazenamento
// 	p.horasRetencao = horasRetencao
// 	p.mu.Unlock()

// 	if err := os.MkdirAll(p.armazenamento, os.ModePerm); err != nil {
// 		return err
// 	}

// 	var nPrcs []Camera
// 	for _, p := range prcsBkp {
// 		nPrcs = append(nPrcs, ))

// 		nPrcs = append(nPrcs, Processo{
// 			ProcessoID:  p.ProcessoID,
// 			EnderecoIP:  p.EnderecoIP,
// 			Porta:       p.Porta,
// 			Canal:       p.Canal,
// 			Usuario:     p.Usuario,
// 			Senha:       p.Senha,
// 			Processador: p.Processador,
// 		})
// 	}

// 	p.StartProcessos(nPrcs)

// 	return nil
// }

// =================================================================================

func (p *Processador) begintHousekeeper() {
	d := time.Now().Add(time.Duration(-p.horasRetencao) * time.Hour)

	err := filepath.Walk(p.armazenamento, func(path string, info os.FileInfo, err error) error {
		if path == p.armazenamento {
			return nil
		}

		if info.ModTime().Before(d) {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		p.errChan <- operrors.OpError{Err: fmt.Errorf("housekeeper stopped: %w", err)}
	}
}
