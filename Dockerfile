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
WORKDIR /bin
ENTRYPOINT ["authService"]