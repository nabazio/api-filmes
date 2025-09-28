package models

import (
	"time"
)

// Filme representa a estrutura completa de um filme
type Filme struct {
	ID              int       `json:"id"`
	Titulo          string    `json:"titulo"`
	Descricao       string    `json:"descricao"`
	AnoLancamento   int       `json:"ano_lancamento"`
	DuracaoMinutos  int       `json:"duracao_minutos"`
	Genero          string    `json:"genero"`
	Diretor         string    `json:"diretor"`
	Avaliacao       float64   `json:"avaliacao"`
	DataCriacao     time.Time `json:"data_criacao"`
	DataAtualizacao time.Time `json:"data_atualizacao"`
}

// FilmeResumo para listagens
type FilmeResumo struct {
	ID            int     `json:"id"`
	Titulo        string  `json:"titulo"`
	AnoLancamento int     `json:"ano_lancamento"`
	Genero        string  `json:"genero"`
	Diretor       string  `json:"diretor"`
	Avaliacao     float64 `json:"avaliacao"`
}

// FilmeParaCriar estrutura para criação (sem ID e timestamps)
type FilmeParaCriar struct {
	Titulo         string   `json:"titulo"`
	Descricao      *string  `json:"descricao,omitempty"`
	AnoLancamento  int      `json:"ano_lancamento"`
	DuracaoMinutos *int     `json:"duracao_minutos,omitempty"`
	Genero         *string  `json:"genero,omitempty"`
	Diretor        *string  `json:"diretor,omitempty"`
	Avaliacao      *float64 `json:"avaliacao,omitempty"`
}

// FilmeParaAtualizar estrutura para atualização (todos campos opcionais)
type FilmeParaAtualizar struct {
	Titulo         *string  `json:"titulo,omitempty"`
	Descricao      *string  `json:"descricao,omitempty"`
	AnoLancamento  *int     `json:"ano_lancamento,omitempty"`
	DuracaoMinutos *int     `json:"duracao_minutos,omitempty"`
	Genero         *string  `json:"genero,omitempty"`
	Diretor        *string  `json:"diretor,omitempty"`
	Avaliacao      *float64 `json:"avaliacao,omitempty"`
}

// Estruturas de resposta
type RespostaFilmes struct {
	Filmes []FilmeResumo `json:"filmes"`
	Total  int           `json:"total"`
}

type RespostaErro struct {
	Erro     string   `json:"erro"`
	Codigo   int      `json:"codigo"`
	Detalhes []string `json:"detalhes,omitempty"`
}

type RespostaSucesso struct {
	Mensagem string      `json:"mensagem"`
	Dados    interface{} `json:"dados,omitempty"`
}
