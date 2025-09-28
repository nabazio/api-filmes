# Makefile
# Comandos simplificados para o projeto

.PHONY: help build run stop clean test docker-build docker-run

# Variáveis
APP_NAME=api-filmes
DOCKER_IMAGE=$(APP_NAME):latest

# Verificar o que está usando a porta
check-port: ## 🔍 Verificar o que está usando a porta 8080
	@echo "$(BLUE)🔍 Verificando porta 8080...$(NC)"
	@lsof -i :8080 || echo "$(GREEN)✅ Porta 8080 livre$(NC)"

# Matar processos que estão usando a porta
kill-port: ## 💀 Matar processo usando porta 8080
	@echo "$(YELLOW)💀 Matando processos na porta 8080...$(NC)"
	@lsof -ti :8080 | xargs kill -9 2>/dev/null || echo "$(GREEN)✅ Nenhum processo encontrado$(NC)"

# Reset completo
reset: ## 🔄 Reset completo (parar tudo e limpar)
	@echo "$(YELLOW)🔄 Reset completo...$(NC)"
	@docker-compose down -v --remove-orphans 2>/dev/null || true
	@pkill -f api-filmes 2>/dev/null || true
	@pkill -f main 2>/dev/null || true
	@lsof -ti :8080 | xargs kill -9 2>/dev/null || true
	@docker system prune -f
	@echo "$(GREEN)✅ Reset concluído$(NC)"

# Atualizar o comando run para fazer verificação automática
run: ## 🚀 Executar aplicação completa
	@echo "$(BLUE)🚀 Iniciando API de Filmes...$(NC)"
	@echo "$(BLUE)🔍 Verificando porta 8080...$(NC)"
	@lsof -ti :8080 | xargs kill -9 2>/dev/null || true
	@docker-compose down 2>/dev/null || true
	@docker-compose up -d
	@echo "$(GREEN)✅ Aplicação iniciada!$(NC)"
	@echo "$(BLUE)🌐 API: http://localhost:8080$(NC)"
	@echo "$(BLUE)🗄️ Adminer: http://localhost:8081$(NC)"

# Comando padrão
help:
	@echo "🎬 API de Filmes - Comandos Disponíveis:"
	@echo ""
	@echo "  make run          - Executar com Docker Compose"
	@echo "  make stop         - Parar todos os containers"
	@echo "  make build        - Build da aplicação Go"
	@echo "  make test         - Executar testes"
	@echo "  make clean        - Limpar containers e volumes"
	@echo "  make logs         - Ver logs da aplicação"
	@echo "  make db-shell     - Conectar no banco via psql"
	@echo "  make docker-build - Build da imagem Docker"
	@echo ""

# Executar aplicação completa
run:
	@echo "🚀 Iniciando API de Filmes..."
	docker-compose up -d
	@echo "✅ Aplicação rodando em http://localhost:8080"
	@echo "🔍 Adminer disponível em http://localhost:8081"

# Parar aplicação
stop:
	@echo "🛑 Parando aplicação..."
	docker-compose down

# Build local da aplicação
build:
	@echo "🔨 Building aplicação..."
	go build -o main ./cmd/server

# Executar testes
test:
	@echo "🧪 Executando testes..."
	go test ./...

# Limpar containers e volumes
clean:
	@echo "🧹 Limpando containers e volumes..."
	docker-compose down -v --remove-orphans
	docker system prune -f

# Ver logs da aplicação
logs:
	docker-compose logs -f api

# Conectar no banco de dados
db-shell:
	docker-compose exec postgres psql -U postgres -d api_filmes

# Build da imagem Docker
docker-build:
	@echo "🐳 Building imagem Docker..."
	docker build -t $(DOCKER_IMAGE) .

# Setup inicial para desenvolvimento
setup:
	@echo "⚙️ Setup inicial do projeto..."
	cp .env.example .env
	@echo "✅ Arquivo .env criado"
	@echo "📝 Edite o arquivo .env conforme necessário"
	@echo "🚀 Execute 'make run' para iniciar"