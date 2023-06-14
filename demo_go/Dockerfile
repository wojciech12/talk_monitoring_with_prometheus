FROM golang:1.20 as builder

WORKDIR /go/src/github.com/wojciech12/order-manager

COPY go.mod go.sum .
RUN go mod download

COPY service/ service/
COPY middleware/ middleware/ 
COPY main.go .

# Statically link all c libraries by disabling cgo
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o main .

FROM ubuntu:22.04

COPY --from=builder /go/src/github.com/wojciech12/order-manager/main .
ENTRYPOINT ["/main"]
