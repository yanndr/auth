FROM golang:1.20 as builder
LABEL authors="yann"

COPY . /go/src
WORKDIR /go/src
run make cert
RUN make build

FROM scratch
COPY --from=builder /go/src/AuthService /bin/AuthService
COPY --from=builder /go/src/config /bin/config
COPY --from=builder /go/src/cert /bin/cert
WORKDIR /bin
ENTRYPOINT ["AuthService"]