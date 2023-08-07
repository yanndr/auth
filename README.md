# Auth Service
Auth service supports account creation and login (username and password) through a gRPC service.
## Run the service

### With Go
Note: locally the service is set up to store the data on a SQLite database. You can change the configuration to use a different database. For now, the service accepts PostgreSQL or SQLite.  

If you have Go 1.20 installed on your computer, you can run the service from the source:
```shell
go run ./cmd/authService/main.go
```
Alternatively, you can first build the program and then run it:

```shell
make build
```
To run it on Linux and MacOS:
```shell
./authService
```
To run it on Windows:
```shell
authservice.exe
```
### With Docker
If you don't have Go installed on your computer, you can run the service with Docker:

Note: with Docker the program is set up to store the data on a PostgreSQL database
```shell
docker compose up authservice
```
This command will start two containers: a PostgreSQL db and the Auth service

## Test the service
### Unit tests:
```shell
make tests
```
Run integration test with Postgresql (need docker)

```shell
make tests-pg
```

### With a simple client
I build a simple client to test the service

build the client:
```shell
    make client
```

The client have 2 sub command create
```shell
 ./authClient --username=test -password=pasw@rd 
```
and auth
```shell
 ./authClient auth --username=test -password=pasw@rd 
```

additional flags are available:
```
-username string
        the username
-password string
        the password
-addr string
        The server address in the format of host:port (default "localhost:50051")
-ca_file string
        The file containing the CA root cert file (default "cert/ca_cert.pem")
-cert_file string
        The file containing the client cert file (default "cert/client_cert.pem")
-key_file string
        The file containing the client key file (default "cert/client_key.pem")
-tls
        Connection uses TLS if true, else plain TCP
  

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
 - make https://www.gnu.org/software/make/
 - sqlc https://sqlc.dev/
 - openssl https://www.openssl.org/
 - mockgen https://github.com/uber-go/mock
 - docker https://www.docker.com/

