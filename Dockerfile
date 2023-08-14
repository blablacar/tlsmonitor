FROM golang:1.21 as builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o ./tlsmonitor

FROM gcr.io/distroless/base

COPY --from=builder /app/tlsmonitor /tlsmonitor

# Health check 
EXPOSE 8080
# Prometheus exporter
EXPOSE 9090


ENTRYPOINT ["/tlsmonitor"]
