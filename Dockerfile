# Dockerfile
# Multi-stage build com Go 1.22

# Estágio 1: Build da aplicação
FROM golang:1.22-alpine AS builder

# Instalar dependências do sistema
RUN apk add --no-cache git ca-certificates tzdata

# Configurar diretório de trabalho
WORKDIR /app

# Configurar variáveis de ambiente
ENV CGO_ENABLED=0 
ENV GOOS=linux 
ENV GOARCH=amd64

# Copiar arquivos de dependências
COPY go.mod go.sum ./

# Download das dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN go build -ldflags='-w -s' -o main ./cmd/server

# Estágio 2: Imagem final
FROM alpine:latest

# Instalar dependências mínimas
RUN apk --no-cache add ca-certificates tzdata netcat-openbsd

# Criar usuário não-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Configurar diretório
WORKDIR /app

# Copiar binário
COPY --from=builder /app/main .

# Copiar scripts
COPY scripts/ ./scripts/
RUN chmod +x ./scripts/*.sh

# Usuário não-root
USER appuser

# Porta
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD nc -z localhost 8080 || exit 1

# Comando
CMD ["./main"]