FROM golang:1.20 as builder
LABEL authors="yann"

COPY . /go/src
WORKDIR /go/src
run make cert
RUN make build

FROM debian:stable-slim

COPY --from=builder /go/src/authService /bin/authService
COPY --from=builder /go/src/config /bin/config
COPY --from=builder /go/src/cert /bin/cert
#RUN apt-get update
#RUN apt-get install -y ca-certificates
#COPY --from=builder go/src/cert/ca_cert.pem /usr/local/share/ca-certificates/ca_cert.pem
#RUN chmod 644 /usr/local/share/ca-certificates/ca_cert.pem && update-ca-certificates
WORKDIR /bin
ENTRYPOINT ["authService"]