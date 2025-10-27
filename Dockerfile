# Etapa de build
FROM golang:1.25.1-alpine AS builder

WORKDIR /app

# Dependencias
RUN apk add --no-cache make bash git

# Copiar go.mod y descargar deps primero para cachear layers
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el proyecto
COPY . .

# Instalar pq
RUN go get github.com/lib/pq@latest

# Compilar
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o my-app .

# Etapa final
FROM alpine:latest

WORKDIR /app

# Copiar binario desde builder
COPY --from=builder /app/my-app . 

# Copiar directorio static
COPY --from=builder /app/static ./static

# Exponer puerto
EXPOSE 8080

# Ejecutar app
CMD ["./my-app"]