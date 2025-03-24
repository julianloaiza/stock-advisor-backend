# ==========================================================
# ETAPA DE COMPILACIÓN
# ==========================================================
FROM golang:1.23-alpine AS builder

# Instalar dependencias del sistema
RUN apk add --no-cache git

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias primero para aprovechar la caché de Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar la aplicación con optimizaciones
RUN CGO_ENABLED=0 GOOS=linux go build -o stock-advisor .

# ==========================================================
# ETAPA FINAL
# ==========================================================
FROM alpine:latest

# Instalar certificados CA para HTTPS y curl para verificaciones de salud
RUN apk --no-cache add ca-certificates curl

# Establecer directorio de trabajo
WORKDIR /app

# Copiar el binario compilado
COPY --from=builder /app/stock-advisor .

# Copiar archivos de configuración necesarios
COPY recommendation_factors.json .

# Exponer el puerto
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./stock-advisor"]