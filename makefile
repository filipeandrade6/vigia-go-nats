.PHONY: run
run:
	docker-compose -f deployments/docker-compose.yml up

.PHONY: test
test:
	# go clean --cache
	go test -cover ./...

.PHONY: create-postgres
create-postgres:
	docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=secret -v ${PWD}/dev/postgresql/:/var/lib/postgresql/data --name dev-postgres postgres

.PHONY: create-pgadmin
create-pgadmin:
	docker run -d -p 80:80 -e 'PGADMIN_DEFAULT_EMAIL=user@domain.local' -e 'PGADMIN_DEFAULT_PASSWORD=secret' --name dev-pgadmin dpage/pgadmin4

.PHONY: postgres
postgres:
	docker restart dev-postgres

.PHONY: pgadmin
pgadmin:
	docker restart dev-pgadmin

.PHONY: docker-clean-volume
docker-clean-volume:
	docker volume ls -qf dangling=true | xargs -r docker volume rm

.PHONY: gerencia
gerencia:
	go run ./cmd/gerencia/main.go --auth-directory=${PWD}/deployments/keys