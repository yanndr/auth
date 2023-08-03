VERSION=$(shell cat cmd/authService/version)
LDFLAGS=-ldflags "-X main.version=${VERSION}"

.PHONY: build
build:
	CGO_ENABLED=0 go build ${LDFLAGS} -o AuthService cmd/authService/main.go

.PHONY: proto
proto:
	protoc api/v1/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.