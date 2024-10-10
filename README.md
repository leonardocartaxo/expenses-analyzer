# expenses-analyzer

Pet project for learning Golang and maybe a boilerplate for future projects

## Run
```shell
go run cmd/api/main.go
```

## Update API docs
```shell
./scripts/gen_doc.sh
```

## Migrations
### Sync DB
```shell
go run cmd/migrate/main.go up
```
Or use auto migrate by setting env var to `DB_AUTO_MIGRATE=true`

### Generate an migration DB
```shell
go run cmd/migrate/main.go create <migration name> sql
```

## **TODO**
* logger
* expenses route
* unit tests
* e2e test
* docker, docker compose
* authentication
* graphql
