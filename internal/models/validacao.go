package models

import (
	"fmt"
	"strings"
	"time"
)

// ValidarFilme valida os dados de um filme antes de salvar
func ValidarFilme(filme *FilmeParaCriar) []string {
	var erros []string

	// Validar título
	if strings.TrimSpace(filme.Titulo) == "" {
		erros = append(erros, "título é obrigatório")
	} else if len(filme.Titulo) > 255 {
		erros = append(erros, "título deve ter no máximo 255 caracteres")
	}

	// Validar ano de lançamento
	anoAtual := time.Now().Year()
	if filme.AnoLancamento < 1888 { // Primeiro filme da história
		erros = append(erros, "ano de lançamento deve ser maior que 1887")
	} else if filme.AnoLancamento > anoAtual+5 { // Máximo 5 anos no futuro
		erros = append(erros, fmt.Sprintf("ano de lançamento não pode ser maior que %d", anoAtual+5))
	}

	// Validar duração
	if filme.DuracaoMinutos != nil && *filme.DuracaoMinutos <= 0 {
		erros = append(erros, "duração deve ser maior que 0 minutos")
	}

	// Validar gênero
	if filme.Genero != nil && len(*filme.Genero) > 100 {
		erros = append(erros, "gênero deve ter no máximo 100 caracteres")
	}

	// Validar diretor
	if filme.Diretor != nil && len(*filme.Diretor) > 255 {
		erros = append(erros, "nome do diretor deve ter no máximo 255 caracteres")
	}

	// Validar avaliação
	if filme.Avaliacao != nil {
		if *filme.Avaliacao < 0 || *filme.Avaliacao > 10 {
			erros = append(erros, "avaliação deve estar entre 0 e 10")
		}
	}

	return erros
}

// ValidarFilmeParaAtualizar valida dados para atualização (campos opcionais)
func ValidarFilmeParaAtualizar(filme *FilmeParaAtualizar) []string {
	var erros []string

	// Validar título (se fornecido)
	if filme.Titulo != nil {
		if strings.TrimSpace(*filme.Titulo) == "" {
			erros = append(erros, "título não pode estar vazio")
		} else if len(*filme.Titulo) > 255 {
			erros = append(erros, "título deve ter no máximo 255 caracteres")
		}
	}

	// Validar ano (se fornecido)
	if filme.AnoLancamento != nil {
		anoAtual := time.Now().Year()
		if *filme.AnoLancamento < 1888 {
			erros = append(erros, "ano de lançamento deve ser maior que 1887")
		} else if *filme.AnoLancamento > anoAtual+5 {
			erros = append(erros, fmt.Sprintf("ano de lançamento não pode ser maior que %d", anoAtual+5))
		}
	}

	// Outras validações similares...
	if filme.DuracaoMinutos != nil && *filme.DuracaoMinutos <= 0 {
		erros = append(erros, "duração deve ser maior que 0 minutos")
	}

	if filme.Avaliacao != nil && (*filme.Avaliacao < 0 || *filme.Avaliacao > 10) {
		erros = append(erros, "avaliação deve estar entre 0 e 10")
	}

	return erros
}
