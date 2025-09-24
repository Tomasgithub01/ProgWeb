# -------------------------
# Etapa de build
# -------------------------
FROM golang:1.25.1-alpine AS builder

# Carpeta de trabajo dentro del contenedor
WORKDIR /app

# Instalar dependencias necesarias
RUN apk add --no-cache make bash git

# Copiar todo el proyecto (incluye código sqlc ya generado)
COPY . .
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN sqlc generate
# Construir la app Go como binario estático
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o my-app .

# -------------------------
# Etapa final: contenedor liviano
# -------------------------
FROM alpine:latest

WORKDIR /app

# Copiar binario desde el builder
COPY --from=builder /app/my-app .

# Copiar archivos estáticos (HTML, etc.)
COPY index.html ./

# Variables de entorno
ENV DB_URL="postgres://admin:#Admin20250915@db:5432/tpespecialweb?sslmode=disable"

# Exponer puerto
EXPOSE 8080

# Ejecutar la app
CMD ["./my-app"]
