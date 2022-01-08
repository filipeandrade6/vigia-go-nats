## VIGIA

## MIGRAR O HOUSEKEEPER PARA O PYTHON


Programa para processamento de metadados de câmeras.

#### JUMPSTART

1. clone o repositório
1. cole os arquivos de aquisição dos metadados
1. cria o banco de dados e interface web (comandos make)
1. crie as tabelas
1. inicie o servidor de gravacao
1. simule servidor de gerencia com o EVANS (está dentro do diretorio dos protofiles)

#### THINKING

* Se der erro por disconexão?
* Novo banco com histórico de notificações
* NewUnit - cria um banco de dados de teste
* NewIntegration - cria um db, alimenta ele  e constroi um autenticador (cria chave, cria um autenticador com essa chave)
Retorna um test { DB, LOG, AUTH, testing.T e função de teardown}
* Token - gera um token autenticado para o usuario
* store - usuarioStore, claims e token utilizando o test acima
* como a verificação de auth fica na requisição da Store - não vou precisar testar o Authentication

#### TODOs

* Comando make para popular banco de dados novo
* Verificar os CASCADE do banco de dados
* Implementar testes
* Trocar nos logs *ERROR* por *error*
* Health check
* Interface no querier
* gRPC em contexto e Health Server [github-1](https://gist.github.com/akhenakh/38dbfea70dc36964e23acc19777f3869) [github-2](https://github.com/grpc/grpc/blob/master/doc/health-checking.md)
* Atualizar armazenamento mover imagens para novo local
* Verificar timezone na aplicação e quando abre o banco de dados
* Colocar interface no querier
* Frontend
* Caso for utilizar servidor e gravacao na mesma máquina, não utilizar protocolo TCP e sim Unix Pipe.

#### TUTORIAIS

* Syslog server, padronização e coleta: [datadog blog](https://www.datadoghq.com/blog/go-logging/#implement-a-standard-logging-interface)

#### Comandos *Makefile*

- `make protobuf` gera os arquivos protobuf de acordo com os arquivos proto em /api/proto/v1
- `make run` executa as aplicações de gerencia e gravação
- `make test` executa os testes
- `make create-postgres` cria o container Docker de Postgres para desenvolvimento
- `make create-pgadmin` cria o container Docker de pgAdmin para desenvolvimento
- `make postgres` reinicia o container Postgres
- `make pgadmin` reinicia o container pgAdmin

### OUTROS

* Fedora
 * Na instalação do protobuf - não instale com dnf install protoc - siga esta [resposta](https://stackoverflow.com/questions/40025602/how-to-use-predifined-protobuf-type-i-e-google-protobuf-timestamp-proto-wit)
