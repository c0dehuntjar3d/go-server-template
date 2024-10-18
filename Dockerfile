FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .
ADD go.sum .

RUN go mod download

COPY . .

WORKDIR /build/cmd

RUN go build -o server main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/cmd/server /build/server
COPY --from=builder /build/docker.env /build/.env

CMD ["/build/server"]
