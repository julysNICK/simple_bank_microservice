DB_URL=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable

postgres:
	docker run --name postgres12 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateupprod:
	migrate -path db/migration -database "postgresql://root:NbgUcBb5kqnC1MbtWbdn@simple-bank.cg6gs0gzgr12.us-east-1.rds.amazonaws.com:5432/simpleBank" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedownprod:
	migrate -path db/migration -database "postgresql://root:NbgUcBb5kqnC1MbtWbdn@simple-bank.cg6gs0gzgr12.us-east-1.rds.amazonaws.com:5432/simpleBank" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1
db_docs:
	dbdocs build doc/db.dbml
db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/julysNICK/simplebank/db/sqlc Store
proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
  --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
  proto/*.proto
evans:
	evans --host localhost --port 9090 -r repl
# history | grep "docker run" sudo systemctl start docker sudo docker start postgres12	

.PHONY: createdb dropdb postgres migrateup migratedown migrateup1 migratedown1 db_docs db_schema sqlc test server mock proto
