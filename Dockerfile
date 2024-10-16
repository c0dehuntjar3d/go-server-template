FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash make gcc gettext musl-dev
COPY ["./go.mod", "./go.sum", "./"]

RUN go mod download

COPY ./ ./
RUN go build -o ./bin/app cmd/main.go

FROM alpine 

COPY --from=builder /usr/local/src/bin/app /

CMD ["/app"]
