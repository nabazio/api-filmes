package config

import (
	"fmt"
	"os"
)

// ConfiguracaoBanco contém as informações de conexão com o banco
type ConfiguracaoBanco struct {
	Host      string
	Porta     string
	Usuario   string
	Senha     string
	NomeBanco string
	SSLMode   string
}

// ObterConfiguracaoBanco retorna a configuração do banco de dados
func ObterConfiguracaoBanco() *ConfiguracaoBanco {
	return &ConfiguracaoBanco{
		Host:      obterVariavelOuPadrao("DB_HOST", "localhost"),
		Porta:     obterVariavelOuPadrao("DB_PORT", "5432"),
		Usuario:   obterVariavelOuPadrao("DB_USER", "postgres"),
		Senha:     obterVariavelOuPadrao("DB_PASSWORD", "160391"),
		NomeBanco: obterVariavelOuPadrao("DB_NAME", "api_filmes"),
		SSLMode:   obterVariavelOuPadrao("DB_SSLMODE", "disable"),
	}
}

// StringConexao gera a string de conexão para o PostgreSQL
func (c *ConfiguracaoBanco) StringConexao() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Porta, c.Usuario, c.Senha, c.NomeBanco, c.SSLMode)
}

// obterVariavelOuPadrao busca uma variável de ambiente ou retorna valor padrão
func obterVariavelOuPadrao(chave, valorPadrao string) string {
	if valor := os.Getenv(chave); valor != "" {
		return valor
	}
	return valorPadrao
}
