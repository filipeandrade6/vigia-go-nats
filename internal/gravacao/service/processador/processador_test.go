package processador_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ardanlabs/service/business/sys/validate"
	"github.com/filipeandrade6/vigia-go/internal/core/registro"
	"github.com/filipeandrade6/vigia-go/internal/core/veiculo"
	"github.com/filipeandrade6/vigia-go/internal/database/tests"
	"github.com/filipeandrade6/vigia-go/internal/gravacao/service/processador"
	"github.com/filipeandrade6/vigia-go/internal/sys/operrors"
)

// TODO teste de adicionar placa
// TODO teste de atualizar matchlist removendo e ver se aparece
// TODO

// TODO colocar campo para a cameras dar diferentes tipos de erro
type CameraTeste struct {
	ProcessoID    string
	Armazenamento string

	closing chan struct{}
}

func (c *CameraTeste) Start(armazenamento string, regChan chan registro.Registro, errChan chan *operrors.OpError) {
	c.closing = make(chan struct{})

	go c.start(
		armazenamento,
		regChan,
		errChan,
	)
}

func (c *CameraTeste) Stop() {
	c.closing <- struct{}{}
	<-c.closing
}

func (c *CameraTeste) GetID() string {
	return c.ProcessoID
}

func (c *CameraTeste) start(armazenamento string, regChan chan registro.Registro, errChan chan *operrors.OpError) {
	defer close(c.closing)

	var i int
	for {
		select {
		case <-c.closing:
			fmt.Println("cancel called")
			return

		default:
			fmt.Print(i, "..")
			time.Sleep(time.Duration(time.Millisecond * 500))
			r := registro.Registro{
				RegistroID:    validate.GenerateID(),
				ProcessoID:    c.ProcessoID,
				Placa:         fmt.Sprintf("ABC%04d", i),
				TipoVeiculo:   "sedan",
				CorVeiculo:    "prata",
				MarcaVeiculo:  "honda",
				Armazenamento: "",
				Confianca:     0.50,
				CriadoEm:      time.Now(),
			}
			r.Armazenamento = fmt.Sprintf("%s/%d_%s", armazenamento, r.CriadoEm.Unix(), r.RegistroID)
			regChan <- r

			err := os.WriteFile(filepath.Join(armazenamento, fmt.Sprintf("%d-%s.txt", i, r.RegistroID)), []byte("hello\n"), 0644)
			if err != nil {
				errChan <- &operrors.OpError{ProcessoID: c.ProcessoID, Err: err}
			}
			i++
		}
	}
}

func TestProcessador(t *testing.T) {
	log, db, teardown := tests.New(t)
	t.Cleanup(teardown)

	registroCore := registro.NewCore(log, db)
	veiculoCore := veiculo.NewCore(log, db)

	ctx := context.Background()

	t.Log("\tTestando Processador...............")
	{
		errChan := make(chan operrors.OpError)
		matchChan := make(chan string)

		// ticker := time.NewTicker(3 * time.Second)

		np := processador.New(
			"/home/filipe/Área de trabalho/tmp/1",
			1,
			registroCore,
			errChan,
			matchChan,
		)

		veiculos, err := veiculoCore.QueryAll(ctx)
		if err != nil {
			t.Fatalf("\t%s\tShould be able to query veiculos: %s.", tests.Failed, err)
		}
		t.Logf("\t%s\tShould be able to query veiculos.", tests.Success)

		var veiculosList []string

		for _, v := range veiculos {
			veiculosList = append(veiculosList, v.Placa)
		}

		np.UpdateMatchlist(veiculosList)

		go np.Start()

		prc := []processador.Camera{&CameraTeste{ProcessoID: "d03307d4-2b28-4c23-a004-3da32e5b8bb1"}}

		np.StartProcessos(prc)

		running, retrying := np.ListProcessos()
		if len(running) != 1 || running[0] != prc[0].GetID() || len(retrying) != 0 {
			t.Fatalf("\t%s\tShould be able to retrieve only the started processo.", tests.Failed)
		}
		t.Logf("\t%s\tShould be able to retrieve only the started processo.", tests.Success)

		// TODO:

	}
}

// 		nProcesso = append(nProcesso, processador.Processo{
// 			ProcessoID:  "d03307d4-2b28-4c23-a004-3da32e5b8a61",
// 			EnderecoIP:  "11.21.31.41",
// 			Porta:       1,
// 			Canal:       1,
// 			Usuario:     "admin",
// 			Senha:       "admin",
// 			Processador: 0,
// 		})

// 		np.StartProcessos(nProcesso)

// 		prcs = np.ListProcessos()
// 		if len(prcs) != 2 {
// 			t.Fatalf("\t%s\tShould be able to retrieve only the started processos.", tests.Failed)
// 		}
// 		t.Logf("\t%s\tShould be able to retrieve only the started processos.", tests.Success)

// 		var registroMatch string
// 		select {
// 		case r := <-matchChan:
// 			registroMatch = r
// 		case <-ticker.C:
// 			t.Fatalf("\t%s\tShould NOT wait more than 5 seconds for match.", tests.Failed)
// 		}

// 		matched, err := registroCore.QueryByID(ctx, registroMatch)
// 		if err != nil {
// 			t.Fatalf("\t%s\tShould be able to retrieve registro by ID: %s.", tests.Failed, err)
// 		}
// 		t.Logf("\t%s\tShould be able to retrieve registro by ID.", tests.Success)

// 		_, err = veiculoCore.QueryByPlaca(ctx, matched.Placa)
// 		if err != nil {
// 			t.Fatalf("\t%s\tShould be able to retrieve placa by registro: %s.", tests.Failed, err)
// 		}
// 		t.Logf("\t%s\tShould be able to retrieve placa by registro.", tests.Success)

// 		status := np.StatusHousekeeper()
// 		if !status {
// 			t.Fatalf("\t%s\tShould be running the housekeeper.", tests.Failed)
// 		}
// 		t.Logf("\t%s\tShould be running the housekeeper.", tests.Success)

// 		np.StopHousekeeper()

// 		status = np.StatusHousekeeper()
// 		if status {
// 			t.Fatalf("\t%s\tShould be stopped the housekeeper.", tests.Failed)
// 		}
// 		t.Logf("\t%s\tShould be stopped the housekeeper.", tests.Success)

// 		np.AtualizarHousekeeper(2)

// 		path, hours := np.GetServidorInfo()
// 		if path != "/home/filipe" || hours != 2 {
// 			t.Fatalf("\t%s\tShould get updated processador info.", tests.Failed)
// 		}
// 		t.Logf("\t%s\tShould be get updated processador info.", tests.Success)

// 		// np.StopGerencia()

// 		select {
// 		case err := <-errChan:
// 			t.Fatalf("\t%s\tShould NOT get any error from channel: %s.", tests.Failed, err)
// 		case <-ticker.C:
// 			t.Logf("\t%s\tShould NOT get any error from channel.", tests.Success)
// 		}

// 		// TODO vai receber erro de already executing
// 	}
// }
