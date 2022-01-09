## VIGIA

Utilizando NATS no lugar de gRPC

#### JUMPSTART

1. clone o repositório
1. cria o banco de dados e interface web (comandos make)
1. crie as tabelas
1. inicie o servidor de gravacao
1. inicie servidor NATS

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
* Atualizar armazenamento mover imagens para novo local
* Verificar timezone na aplicação e quando abre o banco de dados

#### TUTORIAIS

* Syslog server, padronização e coleta: [datadog blog](https://www.datadoghq.com/blog/go-logging/#implement-a-standard-logging-interface)
