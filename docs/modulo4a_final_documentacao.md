# Módulo 4-A: Dockerização Completa da API
## 🐳 Containerização Profissional para Portfolio

### 📖 Objetivos do Módulo
- Dockerizar completamente a aplicação com multi-stage build
- Configurar orquestração com docker-compose
- Implementar scripts de inicialização e health checks
- Resolver problemas comuns de porta e configuração
- Preparar ambiente de desenvolvimento simplificado

---

## 🧠 Conceitos Fundamentais

### O que é Docker e Por que Usar?

Docker é uma plataforma que permite empacotar aplicações em containers isolados. Para nossa API de filmes:

**Antes do Docker:**
```bash
# Para rodar o projeto, alguém precisa:
1. Instalar Go na versão exata
2. Instalar PostgreSQL e configurar
3. Configurar variáveis de ambiente
4. Resolver conflitos de dependências
5. Debugar problemas de "funciona na minha máquina"
```

**Com Docker:**
```bash
# Para rodar o projeto:
make run
# Pronto! ✨
```

### Benefícios para Portfolio
- **Impressiona recrutadores**: Setup em menos de 2 minutos
- **Demonstra conhecimento DevOps**: Containerização é essencial hoje
- **Elimina problemas de ambiente**: Funciona igual em qualquer máquina
- **Facilita demonstrações**: Pode rodar ao vivo em entrevistas

---

## 📁 Estrutura Final do Projeto

```
api-filmes/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   ├── models/
│   ├── database/
│   ├── config/
│   └── validators/
├── scripts/
│   ├── wait-for-db.sh        # Script para aguardar banco
│   └── init-db.sql           # Inicialização do schema
├── .dockerignore             # Otimização do build
├── .env.example              # Template de configuração
├── Dockerfile                # Recipe da aplicação
├── docker-compose.yml        # Orquestração dos serviços
├── Makefile                  # Comandos simplificados
├── README.md                 # Documentação profissional
├── go.mod
└── go.sum
```

---

## 🐳 Dockerfile Otimizado

### Multi-Stage Build Explicado

```dockerfile
# Dockerfile
# Estágio 1: Build da aplicação
FROM golang:1.22-alpine AS builder

# Instalar dependências do sistema
RUN apk add --no-cache git ca-certificates tzdata

# Configurar ambiente de trabalho
WORKDIR /app
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Copiar dependências primeiro (otimização de cache)
COPY go.mod go.sum ./
RUN go mod download

# Copiar código e fazer build
COPY . .
RUN go build -ldflags='-w -s' -o main ./cmd/server

# Estágio 2: Imagem final (runtime)
FROM alpine:latest

# Dependências mínimas para runtime
RUN apk --no-cache add ca-certificates tzdata netcat-openbsd

# Criar usuário não-root (segurança)
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copiar apenas o binário do estágio anterior
COPY --from=builder /app/main .
COPY scripts/ ./scripts/
RUN chmod +x ./scripts/*.sh

# Executar como usuário não-root
USER appuser

EXPOSE 8080

# Health check para monitoramento
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD nc -z localhost 8080 || exit 1

CMD ["./main"]
```

### Por que Multi-Stage?
- **Imagem final pequena**: ~15MB vs ~300MB sem multi-stage
- **Segurança**: Sem ferramentas de desenvolvimento na produção
- **Performance**: Deploy mais rápido com imagens menores

---

## 🔧 Docker Compose para Orquestração

```yaml
# docker-compose.yml
version: '3.8'

services:
  # Banco PostgreSQL
  postgres:
    image: postgres:15-alpine
    container_name: api-filmes-db
    restart: unless-stopped
    environment:
      POSTGRES_DB: api_filmes
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres123
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-db.sql:/docker-entrypoint-initdb.d/init-db.sql:ro
    networks:
      - api-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d api_filmes"]
      interval: 10s
      timeout: 5s
      retries: 5

  # API de Filmes
  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: api-filmes-app
    restart: unless-stopped
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: postgres123
      DB_NAME: api_filmes
      DB_SSLMODE: disable
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - api-network
    command: ["./scripts/wait-for-db.sh", "postgres", "5432", "./main"]

  # Interface de administração do banco
  adminer:
    image: adminer:latest
    container_name: api-filmes-adminer
    restart: unless-stopped
    ports:
      - "8081:8080"
    depends_on:
      - postgres
    networks:
      - api-network

volumes:
  postgres_data:

networks:
  api-network:
    driver: bridge
```

### Características Importantes
- **Health checks**: Garante que banco está pronto antes da API
- **Volumes persistentes**: Dados não são perdidos entre restarts
- **Rede isolada**: Containers comunicam entre si com segurança
- **Dependências**: API só inicia após PostgreSQL estar saudável

---

## 📜 Scripts de Inicialização

### Script de Espera do Banco

```bash
#!/bin/sh
# scripts/wait-for-db.sh
set -e

host="$1"
port="$2"
shift 2
cmd="$@"

echo "🔄 Aguardando banco de dados em $host:$port..."

until nc -z "$host" "$port"; do
  echo "🔄 Banco ainda não está pronto - aguardando..."
  sleep 2
done

echo "✅ Banco de dados está pronto!"
exec $cmd
```

**Por que é necessário?**
- Docker pode iniciar containers rapidamente
- PostgreSQL demora alguns segundos para aceitar conexões
- Script garante que API só inicie após banco estar 100% pronto

### Script de Inicialização do Banco

```sql
-- scripts/init-db.sql
-- Executado automaticamente na primeira inicialização

-- Criar extensões úteis
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Criar tabela principal
CREATE TABLE IF NOT EXISTS filmes (
    id SERIAL PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    descricao TEXT,
    ano_lancamento INTEGER NOT NULL,
    duracao_minutos INTEGER,
    genero VARCHAR(100),
    diretor VARCHAR(255),
    avaliacao DECIMAL(3,1) CHECK (avaliacao >= 0 AND avaliacao <= 10),
    data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_atualizacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para performance
CREATE INDEX IF NOT EXISTS idx_filmes_titulo ON filmes(titulo);
CREATE INDEX IF NOT EXISTS idx_filmes_genero ON filmes(genero);
CREATE INDEX IF NOT EXISTS idx_filmes_ano ON filmes(ano_lancamento);
CREATE INDEX IF NOT EXISTS idx_filmes_avaliacao ON filmes(avaliacao);

-- Dados de exemplo (apenas se tabela vazia)
INSERT INTO filmes (titulo, descricao, ano_lancamento, duracao_minutos, genero, diretor, avaliacao)
SELECT * FROM (VALUES
    ('O Poderoso Chefão', 'A saga de uma família mafiosa italiana nos Estados Unidos', 1972, 175, 'Drama', 'Francis Ford Coppola', 9.2),
    ('Cidade de Deus', 'Retrato da violência urbana no Rio de Janeiro', 2002, 130, 'Drama', 'Fernando Meirelles', 8.6),
    ('Vingadores: Ultimato', 'Os heróis se unem para derrotar Thanos', 2019, 181, 'Ação', 'Anthony e Joe Russo', 8.4),
    ('Parasita', 'Uma família pobre se infiltra na casa de uma família rica', 2019, 132, 'Thriller', 'Bong Joon-ho', 8.6),
    ('Pulp Fiction', 'Histórias entrelaçadas no submundo de Los Angeles', 1994, 154, 'Crime', 'Quentin Tarantino', 8.9)
) AS dados(titulo, descricao, ano_lancamento, duracao_minutos, genero, diretor, avaliacao)
WHERE NOT EXISTS (SELECT 1 FROM filmes LIMIT 1);

-- Trigger para atualizar timestamp automaticamente
CREATE OR REPLACE FUNCTION update_data_atualizacao()
RETURNS TRIGGER AS $$
BEGIN
    NEW.data_atualizacao = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

DROP TRIGGER IF EXISTS trigger_update_data_atualizacao ON filmes;
CREATE TRIGGER trigger_update_data_atualizacao
    BEFORE UPDATE ON filmes
    FOR EACH ROW
    EXECUTE FUNCTION update_data_atualizacao();
```

**Características do script:**
- **Idempotente**: Pode ser executado múltiplas vezes sem problemas
- **Completo**: Inclui schema, índices, dados e triggers
- **Inteligente**: Só insere dados se tabela estiver vazia

---

## ⚙️ Arquivos de Configuração

### .dockerignore para Otimização

```gitignore
# .dockerignore
# Reduz tamanho do contexto de build

# Git e documentação
.git
.gitignore
README.md
docs/
*.md

# Arquivos temporários
*.tmp
*.log
.DS_Store
Thumbs.db

# Dependências locais
vendor/

# Configurações locais
.env
.env.local

# IDEs
.vscode/
.idea/
*.swp
*.swo

# Artifacts de build
main
api-filmes
*.exe

# Docker files (não precisam estar na imagem)
Dockerfile*
docker-compose*
.dockerignore

# CI/CD
.github/
```

### Template de Configuração

```bash
# .env.example
# Copie para .env e ajuste conforme necessário

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres123
DB_NAME=api_filmes
DB_SSLMODE=disable

# Application Configuration
PORT=8080
ENV=development

# Futuras configurações de segurança
# JWT_SECRET=your-secret-key
# API_KEY=your-api-key
```

---

## 🛠️ Makefile Básico

```makefile
# Makefile
.PHONY: help run stop clean logs reset check-port kill-port

# Cores para output visual
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m

APP_NAME=api-filmes

help: ## 📖 Mostrar ajuda
	@echo "$(BLUE)🎬 API de Filmes - Comandos Disponíveis$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

run: ## 🚀 Executar aplicação completa
	@echo "$(BLUE)🚀 Iniciando API de Filmes...$(NC)"
	@echo "$(BLUE)🔍 Verificando porta 8080...$(NC)"
	@lsof -ti :8080 | xargs kill -9 2>/dev/null || true
	@docker-compose down 2>/dev/null || true
	@docker-compose up -d
	@echo "$(GREEN)✅ Aplicação iniciada!$(NC)"
	@echo "$(BLUE)🌐 API: http://localhost:8080$(NC)"
	@echo "$(BLUE)🗄️ Adminer: http://localhost:8081$(NC)"

stop: ## 🛑 Parar aplicação
	@echo "$(YELLOW)🛑 Parando aplicação...$(NC)"
	@docker-compose down
	@echo "$(GREEN)✅ Aplicação parada$(NC)"

logs: ## 📋 Ver logs da aplicação
	@docker-compose logs -f api

clean: ## 🧹 Limpar containers e volumes
	@echo "$(YELLOW)🧹 Limpando containers e volumes...$(NC)"
	@docker-compose down -v --remove-orphans
	@docker system prune -f
	@echo "$(GREEN)✅ Limpeza concluída$(NC)"

reset: ## 🔄 Reset completo (resolver conflitos de porta)
	@echo "$(YELLOW)🔄 Reset completo...$(NC)"
	@docker-compose down -v --remove-orphans 2>/dev/null || true
	@pkill -f api-filmes 2>/dev/null || true
	@pkill -f main 2>/dev/null || true
	@lsof -ti :8080 | xargs kill -9 2>/dev/null || true
	@docker system prune -f
	@echo "$(GREEN)✅ Reset concluído$(NC)"

check-port: ## 🔍 Verificar o que está usando a porta 8080
	@echo "$(BLUE)🔍 Verificando porta 8080...$(NC)"
	@lsof -i :8080 || echo "$(GREEN)✅ Porta 8080 livre$(NC)"

kill-port: ## 💀 Matar processo usando porta 8080
	@echo "$(YELLOW)💀 Matando processos na porta 8080...$(NC)"
	@lsof -ti :8080 | xargs kill -9 2>/dev/null || echo "$(GREEN)✅ Nenhum processo encontrado$(NC)"

setup: ## ⚙️ Setup inicial do projeto
	@echo "$(BLUE)⚙️ Configurando projeto...$(NC)"
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "$(GREEN)✅ Arquivo .env criado$(NC)"; \
	else \
		echo "$(YELLOW)⚠️ Arquivo .env já existe$(NC)"; \
	fi
	@echo "$(GREEN)✅ Setup concluído!$(NC)"
	@echo "$(BLUE)💡 Execute 'make run' para iniciar$(NC)"

dev: setup run ## 🚀 Setup + Run (comando completo)
```

**Comandos mais úteis:**
- `make run` - Inicia tudo automaticamente
- `make reset` - Resolve problemas de porta
- `make logs` - Ver o que está acontecendo
- `make clean` - Limpar tudo para recomeçar

---

## 🧪 Testando a Configuração

### Passo 1: Verificar Arquivos
```bash
# Verificar estrutura
ls -la scripts/
ls -la Dockerfile docker-compose.yml go.mod

# Verificar permissões do script
chmod +x scripts/wait-for-db.sh
```

### Passo 2: Setup Inicial
```bash
# Configurar ambiente
make setup

# Verificar se .env foi criado
cat .env
```

### Passo 3: Primeira Execução
```bash
# Reset completo (caso tenha algo rodando)
make reset

# Iniciar aplicação
make run

# Aguardar ~30 segundos para primeira inicialização
```

### Passo 4: Testar API
```bash
# Verificar saúde da API
curl http://localhost:8080/

# Listar filmes
curl http://localhost:8080/filmes

# Criar novo filme
curl -X POST http://localhost:8080/filmes \
  -H "Content-Type: application/json" \
  -d '{
    "titulo": "Filme Docker",
    "descricao": "Criado via container",
    "ano_lancamento": 2024,
    "duracao_minutos": 120,
    "genero": "Tecnologia",
    "diretor": "DevOps Master",
    "avaliacao": 9.0
  }'

# Verificar estatísticas
curl http://localhost:8080/filmes/estatisticas
```

### Passo 5: Testar Adminer
1. Abra http://localhost:8081
2. Configure conexão:
   - **Sistema**: PostgreSQL
   - **Servidor**: postgres
   - **Usuário**: postgres
   - **Senha**: postgres123
   - **Base de dados**: api_filmes

---

## 🔧 Troubleshooting

### Problema: Porta 8080 em uso
```bash
# Solução rápida
make reset

# Ou verificar e matar manualmente
make check-port
make kill-port
```

### Problema: Build Docker falha
```bash
# Verificar go.mod
cat go.mod
# Deve ter: go 1.21 ou go 1.22 (não 1.25!)

# Corrigir se necessário
go mod tidy

# Rebuild
docker-compose up --build
```

### Problema: Banco não conecta
```bash
# Ver logs do PostgreSQL
docker-compose logs postgres

# Ver logs da API
docker-compose logs api

# Testar conexão manual
docker-compose exec postgres psql -U postgres -d api_filmes -c "SELECT 1;"
```

### Problema: API não responde
```bash
# Verificar se containers estão rodando
docker-compose ps

# Verificar health checks
docker-compose exec api nc -z localhost 8080

# Restart se necessário
make stop && make run
```

### Problema: Dados não persistem
```bash
# Verificar volumes
docker volume ls

# Recrear volumes se necessário
make clean
make run
```

---

## 📊 Métricas de Sucesso

### Performance
- **Build time**: < 60 segundos (primeira vez)
- **Image size**: < 20MB (aplicação final)
- **Startup time**: < 30 segundos (stack completa)
- **Memory usage**: < 50MB por container

### Usabilidade
- **Setup**: 1 comando (`make run`)
- **Dependencies**: Apenas Docker
- **Documentation**: README claro
- **Troubleshooting**: Comandos para resolver problemas

---

## 🎓 Conceitos Aprendidos

### Docker Fundamentals
- **Containerização**: Isolamento de aplicações
- **Multi-stage builds**: Otimização de imagens
- **Volumes**: Persistência de dados
- **Networks**: Comunicação entre containers
- **Health checks**: Monitoramento automático

### DevOps Practices
- **Infrastructure as Code**: docker-compose.yml
- **Automation**: Makefile com comandos padronizados
- **Environment consistency**: Mesmo ambiente em dev/prod
- **Dependency management**: Containers com versões fixas

### Security Best Practices
- **Non-root user**: Containers não executam como root
- **Minimal images**: Alpine Linux para menor superfície de ataque
- **Secrets management**: Preparado para variáveis seguras
- **Network isolation**: Containers em rede isolada

### Observability
- **Health checks**: Docker monitora saúde automaticamente
- **Structured logs**: Logs claros e organizados
- **Service dependencies**: Ordem correta de inicialização
- **Resource monitoring**: Preparado para métricas

---

## ✅ Checklist Final

### Arquivos Criados
- [ ] `Dockerfile` com multi-stage build
- [ ] `docker-compose.yml` completo
- [ ] `scripts/wait-for-db.sh` executável
- [ ] `scripts/init-db.sql` com schema
- [ ] `.dockerignore` otimizado
- [ ] `.env.example` documentado
- [ ] `Makefile` com comandos úteis

### Funcionalidades
- [ ] Containers iniciam na ordem correta
- [ ] API conecta no PostgreSQL
- [ ] Dados persistem entre restarts
- [ ] Health checks funcionam
- [ ] Adminer acessível
- [ ] Comandos make funcionam

### Testes Realizados
- [ ] `make run` inicia tudo sem erros
- [ ] `curl http://localhost:8080/` responde
- [ ] `curl http://localhost:8080/filmes` lista filmes
- [ ] POST de novos filmes funciona
- [ ] Adminer acessa banco corretamente
- [ ] `make reset` resolve conflitos de porta

### Qualidade
- [ ] Build rápido (< 1 minuto)
- [ ] Imagem pequena (< 20MB)
- [ ] Startup rápido (< 30 segundos)
- [ ] Logs claros e úteis
- [ ] Comandos autodocumentados
- [ ] Troubleshooting efetivo

---

## 🎉 Resultado Alcançado

### Transformação Completa
**Antes**: "Clone e torça para funcionar"  
**Depois**: `make run` e está pronto em 2 minutos

### Para seu Portfolio
- **Profissionalismo**: Demonstra conhecimento de containerização
- **Facilidade**: Recrutadores podem testar rapidamente
- **Escalabilidade**: Base para deploy em produção
- **Modernidade**: Uso de tecnologias atuais

### Para Entrevistas
- **Demonstração ao vivo**: Pode rodar durante a entrevista
- **Discussão técnica**: Arquitetura Docker bem fundamentada
- **Problem solving**: Resolução de conflitos de porta
- **DevOps mindset**: Automação e padronização

Sua API agora está completamente containerizada e pronta para impressionar em qualquer portfolio ou entrevista técnica! 🚀