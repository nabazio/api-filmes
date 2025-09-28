package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"api-filmes/internal/database"
	"api-filmes/internal/models"
)

// FilmeHandler contém as dependências para os handlers de filme
type FilmeHandler struct {
	bancoDados *database.BancoDados
}

// NovoFilmeHandler cria uma nova instância do handler
func NovoFilmeHandler(bd *database.BancoDados) *FilmeHandler {
	return &FilmeHandler{bancoDados: bd}
}

// ManipularFilmes lida com requisições para /filmes
func (fh *FilmeHandler) ManipularFilmes(w http.ResponseWriter, r *http.Request) {
	configurarCabecalhos(w)

	switch r.Method {
	case "GET":
		fh.listarFilmes(w, r)
	case "POST":
		fh.criarFilme(w, r)
	default:
		enviarErro(w, "Método não permitido", http.StatusMethodNotAllowed, nil)
	}
}

// ManipularFilmeIndividual lida com requisições para /filmes/{id}
func (fh *FilmeHandler) ManipularFilmeIndividual(w http.ResponseWriter, r *http.Request) {
	configurarCabecalhos(w)

	// Extrair ID da URL
	caminho := strings.TrimPrefix(r.URL.Path, "/filmes/")
	if caminho == "" {
		enviarErro(w, "ID do filme é obrigatório", http.StatusBadRequest, nil)
		return
	}

	id, err := strconv.Atoi(caminho)
	if err != nil {
		enviarErro(w, "ID inválido", http.StatusBadRequest, []string{"ID deve ser um número inteiro"})
		return
	}

	switch r.Method {
	case "GET":
		fh.buscarFilmePorID(w, r, id)
	case "PUT":
		fh.atualizarFilme(w, r, id)
	case "DELETE":
		fh.deletarFilme(w, r, id)
	default:
		enviarErro(w, "Método não permitido", http.StatusMethodNotAllowed, nil)
	}
}

// listarFilmes retorna todos os filmes
func (fh *FilmeHandler) listarFilmes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("📋 Listando filmes...")

	filmes, err := fh.bancoDados.BuscarTodosFilmes()
	if err != nil {
		fmt.Printf("❌ Erro ao buscar filmes: %v\n", err)
		enviarErro(w, "Erro interno do servidor", http.StatusInternalServerError, nil)
		return
	}

	total, err := fh.bancoDados.ContarFilmes()
	if err != nil {
		fmt.Printf("⚠️ Erro ao contar filmes: %v\n", err)
		total = len(filmes)
	}

	resposta := models.RespostaFilmes{
		Filmes: filmes,
		Total:  total,
	}

	fmt.Printf("✅ Listados %d filmes\n", len(filmes))
	enviarJSON(w, resposta, http.StatusOK)
}

// buscarFilmePorID retorna um filme específico
func (fh *FilmeHandler) buscarFilmePorID(w http.ResponseWriter, r *http.Request, id int) {
	fmt.Printf("🔍 Buscando filme ID: %d\n", id)

	filme, err := fh.bancoDados.BuscarFilmePorID(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			enviarErro(w, fmt.Sprintf("Filme com ID %d não encontrado", id), http.StatusNotFound, nil)
		} else {
			fmt.Printf("❌ Erro ao buscar filme: %v\n", err)
			enviarErro(w, "Erro interno do servidor", http.StatusInternalServerError, nil)
		}
		return
	}

	fmt.Printf("✅ Filme encontrado: %s\n", filme.Titulo)
	enviarJSON(w, filme, http.StatusOK)
}

// criarFilme cria um novo filme
func (fh *FilmeHandler) criarFilme(w http.ResponseWriter, r *http.Request) {
	fmt.Println("➕ Criando novo filme...")

	var filme models.FilmeParaCriar

	// Decodificar JSON do body
	if err := json.NewDecoder(r.Body).Decode(&filme); err != nil {
		enviarErro(w, "JSON inválido", http.StatusBadRequest, []string{"Verifique a sintaxe do JSON"})
		return
	}

	// Validar dados
	if erros := models.ValidarFilme(&filme); len(erros) > 0 {
		enviarErro(w, "Dados inválidos", http.StatusBadRequest, erros)
		return
	}

	// Salvar no banco
	novoFilme, err := fh.bancoDados.CriarFilme(&filme)
	if err != nil {
		fmt.Printf("❌ Erro ao criar filme: %v\n", err)
		enviarErro(w, "Erro interno do servidor", http.StatusInternalServerError, nil)
		return
	}

	fmt.Printf("✅ Filme criado: %s (ID: %d)\n", novoFilme.Titulo, novoFilme.ID)

	resposta := models.RespostaSucesso{
		Mensagem: "Filme criado com sucesso",
		Dados:    novoFilme,
	}

	enviarJSON(w, resposta, http.StatusCreated)
}

// atualizarFilme atualiza um filme existente
func (fh *FilmeHandler) atualizarFilme(w http.ResponseWriter, r *http.Request, id int) {
	fmt.Printf("✏️ Atualizando filme ID: %d\n", id)

	var filme models.FilmeParaAtualizar

	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&filme); err != nil {
		enviarErro(w, "JSON inválido", http.StatusBadRequest, []string{"Verifique a sintaxe do JSON"})
		return
	}

	// Validar dados
	if erros := models.ValidarFilmeParaAtualizar(&filme); len(erros) > 0 {
		enviarErro(w, "Dados inválidos", http.StatusBadRequest, erros)
		return
	}

	// Atualizar no banco
	filmeAtualizado, err := fh.bancoDados.AtualizarFilme(id, &filme)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			enviarErro(w, fmt.Sprintf("Filme com ID %d não encontrado", id), http.StatusNotFound, nil)
		} else {
			fmt.Printf("❌ Erro ao atualizar filme: %v\n", err)
			enviarErro(w, "Erro interno do servidor", http.StatusInternalServerError, nil)
		}
		return
	}

	fmt.Printf("✅ Filme atualizado: %s\n", filmeAtualizado.Titulo)

	resposta := models.RespostaSucesso{
		Mensagem: "Filme atualizado com sucesso",
		Dados:    filmeAtualizado,
	}

	enviarJSON(w, resposta, http.StatusOK)
}

// deletarFilme remove um filme
func (fh *FilmeHandler) deletarFilme(w http.ResponseWriter, r *http.Request, id int) {
	fmt.Printf("🗑️ Deletando filme ID: %d\n", id)

	err := fh.bancoDados.DeletarFilme(id)
	if err != nil {
		if strings.Contains(err.Error(), "não encontrado") {
			enviarErro(w, fmt.Sprintf("Filme com ID %d não encontrado", id), http.StatusNotFound, nil)
		} else {
			fmt.Printf("❌ Erro ao deletar filme: %v\n", err)
			enviarErro(w, "Erro interno do servidor", http.StatusInternalServerError, nil)
		}
		return
	}

	fmt.Printf("✅ Filme deletado (ID: %d)\n", id)

	resposta := models.RespostaSucesso{
		Mensagem: "Filme deletado com sucesso",
	}

	enviarJSON(w, resposta, http.StatusOK)
}

// Funções utilitárias
func configurarCabecalhos(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func enviarJSON(w http.ResponseWriter, dados interface{}, status int) {
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(dados); err != nil {
		fmt.Printf("❌ Erro ao codificar JSON: %v\n", err)
		http.Error(w, "Erro interno", http.StatusInternalServerError)
	}
}

func enviarErro(w http.ResponseWriter, mensagem string, status int, detalhes []string) {
	erro := models.RespostaErro{
		Erro:     mensagem,
		Codigo:   status,
		Detalhes: detalhes,
	}

	enviarJSON(w, erro, status)
}
