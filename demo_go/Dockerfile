FROM golang:1.10  as builder

WORKDIR /go/src/github.com/wojciech12/order-manager

COPY service/ service/
COPY middleware/ middleware/ 
COPY main.go .

RUN go get -d -v ./...

# Statically link all c libraries by disabling cgo
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o main .

FROM phusion/baseimage
COPY --from=builder /go/src/github.com/wojciech12/order-manager/main .
ENTRYPOINT ["/main"]
