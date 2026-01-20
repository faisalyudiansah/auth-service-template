include ${CURDIR}/configs/envs/.env

run:
	@go run ${CURDIR}/cmd/api/main.go serve-all

nodemon:
	nodemon --exec go run ${CURDIR}/cmd/api/main.go --signal SIGTERM

test:
	@go test -v -cover ${CURDIR}/... 

test-cover:
	@go test ${CURDIR}/... -coverprofile=coverage.out
	@go tool cover -html=coverage.out && rm -f coverage.out

coverage:
	@go test ${CURDIR}/... -coverprofile=cover.out
	@go tool cover -html=cover.out && rm -rf cover.out

mockery:
	@mockery --all --case underscore --dir ${input} --output ${output}

build:
	@go build -o ${CURDIR}/build/main ${CURDIR}/cmd/api/main.go

migratecreate:
	@migrate create -ext sql -dir ${CURDIR}/db/migrations/ -seq ${name}

migrateforce:
	@migrate -path ${CURDIR}/db/migrations/ -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -verbose force 1

migratedown:
	@migrate -path ${CURDIR}/db/migrations/ -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -verbose down

migrateup:
	@migrate -path ${CURDIR}/db/migrations/ -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -verbose up

seed:
	@${CURDIR}/scripts/seed.sh ${DB_USER} ${DB_PASSWORD} ${DB_HOST} ${DB_PORT} ${DB_NAME} ${DB_SSL_MODE}

reset-db:
	@echo "\n migrate down"
	@yes | migrate -path ${CURDIR}/db/migrations/ \
	-database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" \
	-verbose down

	@echo "\n migrate up"
	@migrate -path ${CURDIR}/db/migrations/ \
	-database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" \
	-verbose up

	@echo "\n seeding"
	@${CURDIR}/scripts/seed.sh ${DB_USER} ${DB_PASSWORD} ${DB_HOST} ${DB_PORT} ${DB_NAME} ${DB_SSL_MODE}

v3mock:
	mockery