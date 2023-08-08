VERSION=$(shell cat cmd/authService/version)
LDFLAGS=-ldflags "-X main.Version=${VERSION}"


build:
	go build ${LDFLAGS} cmd/authService/authService.go

client:
	go build cmd/client/client.go

.PHONY: tests
tests:
	go test ./pkg/... -race

.PHONY: proto
proto:
	protoc --go_out=. \
		--go-grpc_out=. \
		--proto_path=. \
		proto/auth.proto

.PHONY: sql
sql:
	sqlc generate

.PHONY: cert
cert:
	cd cert && ./cert.sh

mocks:
	mockgen -source=./pkg/stores/store.go -destination=./pkg/tests/mockStore.go -package=tests
	mockgen -source=./pkg/jwt/jwt.go -destination=./pkg/tests/mockJwt.go -package=tests
	mockgen -source=./pkg/validators/validators.go -destination=./pkg/tests/mockValidators.go -package=tests
	mockgen -source=./pkg/services/userService.go -destination=./pkg/tests/mockUserService.go -package=tests
	mockgen -source=./pkg/services/authService.go -destination=./pkg/tests/mockAuthService.go -package=tests

docker-service:
	docker build -t auth_authservice:latest .
	docker build -t auth_authservice:${VERSION} .

tests-pg:
	docker build -t auth_db -f Dockerfile_postgres .
	docker run  -d --rm --name auth_db_test \
		--env POSTGRES_PASSWORD=passw@rd \
		--env AUTH_USER_PWD=autPassw@ord \
		--env AUTH_DB=auth \
		--env AUTH_USER=auth_user\
		-p 5433:5432 \
		 auth_db
	sleep 5
	go test -tags=pg_test ./pkg/... -race
	docker stop auth_db_test