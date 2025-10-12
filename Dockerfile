
# Etapa de build
FROM golang:1.25.1-alpine AS builder

# Carpeta raiz dentro del contenedor
WORKDIR /app

# Dependencias
RUN apk add --no-cache make bash git

# Copiar el protecto 
COPY . .
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

RUN go get github.com/lib/pq@latest


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o my-app .

RUN sqlc generate

COPY go.mod go.sum ./
RUN go mod download

COPY . .


# Contenedor
FROM alpine:latest

WORKDIR /app

# Copiar binario desde el builder
COPY --from=builder /app/my-app .

# Copiar archivos estáticos (HTML)
COPY index.html ./

# Variables de entorno
ENV DB_URL="postgres://admin:#Admin20250915@db:5432/tpespecialweb?sslmode=disable"

# Exponer puerto
EXPOSE 8080

# Ejecutar la app
CMD ["./my-app"]
