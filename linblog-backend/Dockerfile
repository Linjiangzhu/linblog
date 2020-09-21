FROM golang:latest AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/linblog-backend
COPY . .

RUN go get -d -v ./...
RUN go build -o bin/linblog-backend cmd/server/main.go

FROM alpine

WORKDIR /root/
COPY --from=builder /go/src/linblog-backend/bin/linblog-backend .
COPY --from=builder /go/src/linblog-backend/config.yaml ./conf/config.yaml

CMD ["./linblog-backend", "serve"]

EXPOSE 9090