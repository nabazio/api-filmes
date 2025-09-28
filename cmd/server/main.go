package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"api-filmes/internal/database"
	"api-filmes/internal/handlers"
)

func main() {
	fmt.Println("🎬 Servidor da API de Filmes iniciando...")

	// Conectar ao banco
	bancoDados, err := database.NovaConexao()
	if err != nil {
		log.Fatal("❌ Erro ao conectar com banco:", err)
	}

	// Garantir que a conexão seja fechada ao final
	defer func() {
		if err := bancoDados.Fechar(); err != nil {
			log.Printf("⚠️ Erro ao fechar conexão: %v", err)
		} else {
			fmt.Println("🔌 Conexão com banco fechada")
		}
	}()

	// Criar handler de filmes
	filmeHandler := handlers.NovoFilmeHandler(bancoDados)

	// Configurar rotas com middleware de log
	http.HandleFunc("/", handlers.LogMiddleware(paginaInicial))
	http.HandleFunc("/filmes", handlers.LogMiddleware(filmeHandler.ManipularFilmes))
	http.HandleFunc("/filmes/", handlers.LogMiddleware(filmeHandler.ManipularFilmeIndividual))

	// Adicionar rota para health check
	http.HandleFunc("/health", handlers.LogMiddleware(healthCheck))

	// Iniciar servidor
	porta := ":8080"
	fmt.Printf("🚀 Servidor rodando em http://localhost%s\n", porta)
	fmt.Println("📋 Endpoints disponíveis:")
	fmt.Println("   GET    /              - Informações da API")
	fmt.Println("   GET    /health        - Status do sistema")
	fmt.Println("   GET    /filmes        - Listar todos os filmes")
	fmt.Println("   POST   /filmes        - Criar novo filme")
	fmt.Println("   GET    /filmes/{id}   - Buscar filme por ID")
	fmt.Println("   PUT    /filmes/{id}   - Atualizar filme")
	fmt.Println("   DELETE /filmes/{id}   - Deletar filme")

	if err := http.ListenAndServe(porta, nil); err != nil {
		log.Fatal("❌ Erro ao iniciar servidor:", err)
	}
}

// Página inicial com informações da API
func paginaInicial(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resposta := map[string]interface{}{
		"mensagem": "🎬 Bem-vindo à API de Filmes!",
		"versao":   "2.0.0",
		"recursos": map[string][]string{
			"filmes": {
				"GET /filmes - Lista todos os filmes",
				"POST /filmes - Cria novo filme",
				"GET /filmes/{id} - Busca filme por ID",
				"PUT /filmes/{id} - Atualiza filme",
				"DELETE /filmes/{id} - Remove filme",
			},
			"sistema": {
				"GET /health - Status do sistema",
			},
		},
		"exemplo_criacao": map[string]interface{}{
			"titulo":          "Nome do Filme",
			"descricao":       "Descrição do filme",
			"ano_lancamento":  2024,
			"duracao_minutos": 120,
			"genero":          "Drama",
			"diretor":         "Nome do Diretor",
			"avaliacao":       8.5,
		},
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resposta)
}

// Health check para monitoramento
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resposta := map[string]interface{}{
		"status":    "OK",
		"timestamp": time.Now().Format(time.RFC3339),
		"servico":   "API de Filmes",
		"versao":    "2.0.0",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resposta)
}
