# Etapa de construcción
FROM golang:1.23 AS builder

WORKDIR /builder

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payroll-checker-backend /builder/cmd/main.go

# Etapa de ejecución
FROM alpine:latest

COPY --from=builder /builder/payroll-checker-backend /payroll-checker-backend

EXPOSE 8080

ENTRYPOINT ["/payroll-checker-backend"]
