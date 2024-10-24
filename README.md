# expenses-analyzer

Pet project for learning Golang and maybe a boilerplate for future projects

## Run
```shell
go run cmd/api/main.go
```

## Docs
http://localhost:8080/swagger
### Update API docs
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
Enable DB_AUTO_MIGRATE and DB_DEBUG flags to see on the migration logs query and then create a migration with this command:
```shell
go run cmd/migrate/main.go create <migration name> sql
```

## TODO
* fix api response for duplicated items on DB
* expand expenses module
* unit tests
* e2e test
* docker, docker compose
* authentication
* graphql
