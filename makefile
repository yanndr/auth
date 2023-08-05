VERSION=$(shell cat cmd/authService/version)
LDFLAGS=-ldflags "-X main.version=${VERSION}"


build:
	go build ${LDFLAGS} -o AuthService cmd/authService/main.go

.PHONY: tests
tests:
	go test ./... -race

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
	mockgen -source=./pkg/store/store.go -destination=./pkg/tests/mockStore.go -package=tests
	mockgen -source=./pkg/jwt/jwt.go -destination=./pkg/tests/mockJwt.go -package=tests
	mockgen -source=./pkg/validators/validators.go -destination=./pkg/tests/mockValidators.go -package=tests