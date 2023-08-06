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
	mockgen -source=./pkg/stores/store.go -destination=./pkg/tests/mockStore.go -package=tests
	mockgen -source=./pkg/jwt/jwt.go -destination=./pkg/tests/mockJwt.go -package=tests
	mockgen -source=./pkg/validators/validators.go -destination=./pkg/tests/mockValidators.go -package=tests
	#mockgen -source=./pkg/pb/auth_grpc.pb.go -destination=./pkg/tests/mockAuth_grpc.go -package=tests
	mockgen -source=./pkg/services/userService.go -destination=./pkg/tests/mockUserService.go -package=tests
	mockgen -source=./pkg/services/authentication.go -destination=./pkg/tests/mockAuthentication.go -package=tests

docker-service:
	docker build -t auth_authservice:latest .
	docker build -t auth_authservice:${VERSION} .