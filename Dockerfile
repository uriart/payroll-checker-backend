# Etapa de construcción
FROM golang:1.23 AS builder

WORKDIR /app

COPY app .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payroll-checker-backend main.go

# Etapa de ejecución
FROM alpine:latest

RUN mkdir /renders

COPY --from=builder /app/payroll-checker-backend /payroll-checker-backend
COPY --from=builder /app/payroll_sample.pdf /payroll_sample.pdf

EXPOSE 8080

ENTRYPOINT ["/payroll-checker-backend"]
