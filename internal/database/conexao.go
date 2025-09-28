package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"api-filmes/internal/config"
	"api-filmes/internal/models"

	_ "github.com/lib/pq"
)

type BancoDados struct {
	conexao *sql.DB
}

func NovaConexao() (*BancoDados, error) {
	configuracao := config.ObterConfiguracaoBanco()

	fmt.Println("🔌 Conectando ao banco de dados...")
	fmt.Printf("📍 Host: %s:%s | Banco: %s\n",
		configuracao.Host, configuracao.Porta, configuracao.NomeBanco)

	conexao, err := sql.Open("postgres", configuracao.StringConexao())
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão: %v", err)
	}

	if err := conexao.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao conectar com banco: %v", err)
	}

	fmt.Println("✅ Conexão com banco estabelecida com sucesso!")
	return &BancoDados{conexao: conexao}, nil
}

func (bd *BancoDados) Fechar() error {
	if bd.conexao != nil {
		return bd.conexao.Close()
	}
	return nil
}

// Operações de leitura (já existentes)
func (bd *BancoDados) BuscarTodosFilmes() ([]models.FilmeResumo, error) {
	query := `
        SELECT id, titulo, ano_lancamento, genero, diretor, avaliacao 
        FROM filmes 
        ORDER BY id ASC
    `

	linhas, err := bd.conexao.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query: %v", err)
	}
	defer linhas.Close()

	var filmes []models.FilmeResumo

	for linhas.Next() {
		var filme models.FilmeResumo

		err := linhas.Scan(
			&filme.ID,
			&filme.Titulo,
			&filme.AnoLancamento,
			&filme.Genero,
			&filme.Diretor,
			&filme.Avaliacao,
		)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler dados do filme: %v", err)
		}

		filmes = append(filmes, filme)
	}

	if err := linhas.Err(); err != nil {
		return nil, fmt.Errorf("erro durante leitura dos resultados: %v", err)
	}

	return filmes, nil
}

func (bd *BancoDados) BuscarFilmePorID(id int) (*models.Filme, error) {
	query := `
        SELECT id, titulo, descricao, ano_lancamento, duracao_minutos, 
               genero, diretor, avaliacao, data_criacao, data_atualizacao
        FROM filmes 
        WHERE id = $1
    `

	var filme models.Filme

	err := bd.conexao.QueryRow(query, id).Scan(
		&filme.ID,
		&filme.Titulo,
		&filme.Descricao,
		&filme.AnoLancamento,
		&filme.DuracaoMinutos,
		&filme.Genero,
		&filme.Diretor,
		&filme.Avaliacao,
		&filme.DataCriacao,
		&filme.DataAtualizacao,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("filme com ID %d não encontrado", id)
		}
		return nil, fmt.Errorf("erro ao buscar filme: %v", err)
	}

	return &filme, nil
}

func (bd *BancoDados) ContarFilmes() (int, error) {
	var total int
	query := "SELECT COUNT(*) FROM filmes"
	err := bd.conexao.QueryRow(query).Scan(&total)

	if err != nil {
		return 0, fmt.Errorf("erro ao contar filmes: %v", err)
	}

	return total, nil
}

// ✨ NOVAS OPERAÇÕES DE ESCRITA

// CriarFilme insere um novo filme no banco
func (bd *BancoDados) CriarFilme(filme *models.FilmeParaCriar) (*models.Filme, error) {
	query := `
        INSERT INTO filmes (titulo, descricao, ano_lancamento, duracao_minutos, genero, diretor, avaliacao)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, data_criacao, data_atualizacao
    `

	var novoFilme models.Filme

	// Mapear dados de entrada para filme completo
	novoFilme.Titulo = filme.Titulo
	novoFilme.AnoLancamento = filme.AnoLancamento

	// Tratar campos opcionais
	var descricao, genero, diretor *string
	var duracaoMinutos *int
	var avaliacao *float64

	if filme.Descricao != nil {
		descricao = filme.Descricao
		novoFilme.Descricao = *filme.Descricao
	}
	if filme.Genero != nil {
		genero = filme.Genero
		novoFilme.Genero = *filme.Genero
	}
	if filme.Diretor != nil {
		diretor = filme.Diretor
		novoFilme.Diretor = *filme.Diretor
	}
	if filme.DuracaoMinutos != nil {
		duracaoMinutos = filme.DuracaoMinutos
		novoFilme.DuracaoMinutos = *filme.DuracaoMinutos
	}
	if filme.Avaliacao != nil {
		avaliacao = filme.Avaliacao
		novoFilme.Avaliacao = *filme.Avaliacao
	}

	// Executar inserção
	err := bd.conexao.QueryRow(
		query,
		filme.Titulo,
		descricao,
		filme.AnoLancamento,
		duracaoMinutos,
		genero,
		diretor,
		avaliacao,
	).Scan(&novoFilme.ID, &novoFilme.DataCriacao, &novoFilme.DataAtualizacao)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar filme: %v", err)
	}

	return &novoFilme, nil
}

// AtualizarFilme atualiza um filme existente
func (bd *BancoDados) AtualizarFilme(id int, filme *models.FilmeParaAtualizar) (*models.Filme, error) {
	// Primeiro, verificar se filme existe
	filmeExistente, err := bd.BuscarFilmePorID(id)
	if err != nil {
		return nil, err
	}

	// Construir query dinâmica baseada nos campos fornecidos
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if filme.Titulo != nil {
		setParts = append(setParts, fmt.Sprintf("titulo = $%d", argIndex))
		args = append(args, *filme.Titulo)
		argIndex++
	}

	if filme.Descricao != nil {
		setParts = append(setParts, fmt.Sprintf("descricao = $%d", argIndex))
		args = append(args, *filme.Descricao)
		argIndex++
	}

	if filme.AnoLancamento != nil {
		setParts = append(setParts, fmt.Sprintf("ano_lancamento = $%d", argIndex))
		args = append(args, *filme.AnoLancamento)
		argIndex++
	}

	if filme.DuracaoMinutos != nil {
		setParts = append(setParts, fmt.Sprintf("duracao_minutos = $%d", argIndex))
		args = append(args, *filme.DuracaoMinutos)
		argIndex++
	}

	if filme.Genero != nil {
		setParts = append(setParts, fmt.Sprintf("genero = $%d", argIndex))
		args = append(args, *filme.Genero)
		argIndex++
	}

	if filme.Diretor != nil {
		setParts = append(setParts, fmt.Sprintf("diretor = $%d", argIndex))
		args = append(args, *filme.Diretor)
		argIndex++
	}

	if filme.Avaliacao != nil {
		setParts = append(setParts, fmt.Sprintf("avaliacao = $%d", argIndex))
		args = append(args, *filme.Avaliacao)
		argIndex++
	}

	// Se nenhum campo foi fornecido, retornar filme existente
	if len(setParts) == 0 {
		return filmeExistente, nil
	}

	// Adicionar atualização automática do timestamp
	setParts = append(setParts, fmt.Sprintf("data_atualizacao = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Adicionar WHERE clause
	args = append(args, id)

	query := fmt.Sprintf(`
        UPDATE filmes 
        SET %s 
        WHERE id = $%d
    `, strings.Join(setParts, ", "), argIndex)

	_, err = bd.conexao.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro ao atualizar filme: %v", err)
	}

	// Retornar filme atualizado
	return bd.BuscarFilmePorID(id)
}

// DeletarFilme remove um filme do banco
func (bd *BancoDados) DeletarFilme(id int) error {
	// Primeiro verificar se filme existe
	_, err := bd.BuscarFilmePorID(id)
	if err != nil {
		return err // Já retorna erro adequado (não encontrado ou erro de banco)
	}

	query := "DELETE FROM filmes WHERE id = $1"

	result, err := bd.conexao.Exec(query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar filme: %v", err)
	}

	// Verificar se alguma linha foi afetada
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar deleção: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("filme com ID %d não foi deletado", id)
	}

	return nil
}
