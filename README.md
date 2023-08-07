# Auth Service

## Run the service

### with Go
note: locally the service is set up to store the data on a sqlLite db. You can change the configuration to use a different db. for now service accept postgres or sqlite  

If you have Go 1.20 installed on your computer you can run the service from the source:
```shell
go run ./cmd/authService/main.go
```
alternatively you can first build the program and then run it:

```shell
make build
```
or
```shell
go build ./cmd/authService/
```
Linux and mac:
```shell
./authService
```
Windows:
```shell
authservice.exe
```
### With Docker
If you don't have go install on your computer, you can run the service with docker:

Note with docker the program is set up to store the data on a Postgresql database
```shell
docker compose up
```
This command will start two container: a postgres db and the Auth service

## Test the service
run the unit test:
```shell
make tests
```
Run integration test with Postgresql (need docker)

```shell
make tests-pg
```

## Notes
### TLS mode
If you want to test the service with TLS enabled do the following:
```shell
make cert
```
then update the config file in config/config.yml and change the value TLSConfig useTLS to true:  
```yaml
TLSConfig:
  useTLS: 'true'
```
### Tool used
 - make
 - sqlc
 - openssl
 - mockgen
 - protoc

