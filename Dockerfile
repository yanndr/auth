FROM golang:1.20 as builder
LABEL authors="yann"

COPY . /go/src
WORKDIR /go/src

RUN make build

FROM scratch
COPY --from=builder /go/src/AuthService /bin/AuthService
COPY --from=builder /go/src/config /bin/config
WORKDIR /bin
ENTRYPOINT ["AuthService"]