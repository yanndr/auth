VERSION=$(shell cat cmd/authService/version)
LDFLAGS=-ldflags "-X main.version=${VERSION}"


build:
	CGO_ENABLED=0 go build ${LDFLAGS} -o AuthService cmd/authService/main.go

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