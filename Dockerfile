FROM golang:1.17-alpine3.14 as builder
WORKDIR go/src/k8s-leader-election
COPY go.mod go.sum ./
RUN go mod download
COPY main.go .
RUN go build -v -o /go/bin/leader-election main.go

FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/leader-election /usr/local/bin
USER 65534
EXPOSE 8080
ENTRYPOINT [ "leader-election" ]
