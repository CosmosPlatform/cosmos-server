FROM golang:1.23.1 AS builder
LABEL intermediateStageToBeDeleted=true

RUN mkdir -p /build
WORKDIR /build/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/cosmos-server .

# Stage 2
FROM alpine:3.20

COPY --from=builder /build/bin/cosmos-server /cosmos-server

ENTRYPOINT ["/cosmos-server"]