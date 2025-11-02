FROM golang:alpine3.22 AS builder
LABEL intermediateStageToBeDeleted=true

RUN mkdir -p /build
WORKDIR /build/

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/cosmos-server .

# Stage 2
FROM alpine:3.22.2

COPY --from=builder /build/bin/cosmos-server /cosmos-server

ENTRYPOINT ["/cosmos-server"]