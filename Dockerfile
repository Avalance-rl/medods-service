FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY cmd ./cmd
COPY internal ./internal
COPY go.mod go.sum ./

RUN export GOPROXY=direct
RUN go mod download

WORKDIR /build/cmd/app

RUN go build -o /build/medods-service .

FROM alpine:latest
WORKDIR /app
COPY configs/main.yml ./configs/
COPY --from=builder /build/medods-service /app/medods-service

CMD ["/app/medods-service"]
