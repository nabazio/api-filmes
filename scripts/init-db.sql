-- scripts/init-db.sql
-- Script de inicialização do banco de dados

-- Criar extensões úteis
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Criar tabela de filmes (caso não exista)
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

-- Criar índices para performance
CREATE INDEX IF NOT EXISTS idx_filmes_titulo ON filmes(titulo);
CREATE INDEX IF NOT EXISTS idx_filmes_genero ON filmes(genero);
CREATE INDEX IF NOT EXISTS idx_filmes_ano ON filmes(ano_lancamento);
CREATE INDEX IF NOT EXISTS idx_filmes_avaliacao ON filmes(avaliacao);

-- Inserir dados de exemplo (apenas se tabela estiver vazia)
INSERT INTO filmes (titulo, descricao, ano_lancamento, duracao_minutos, genero, diretor, avaliacao)
SELECT * FROM (VALUES
    ('O Poderoso Chefão', 'A saga de uma família mafiosa italiana nos Estados Unidos', 1972, 175, 'Drama', 'Francis Ford Coppola', 9.2),
    ('Cidade de Deus', 'Retrato da violência urbana no Rio de Janeiro', 2002, 130, 'Drama', 'Fernando Meirelles', 8.6),
    ('Vingadores: Ultimato', 'Os heróis se unem para derrotar Thanos', 2019, 181, 'Ação', 'Anthony e Joe Russo', 8.4),
    ('Parasita', 'Uma família pobre se infiltra na casa de uma família rica', 2019, 132, 'Thriller', 'Bong Joon-ho', 8.6),
    ('Pulp Fiction', 'Histórias entrelaçadas no submundo de Los Angeles', 1994, 154, 'Crime', 'Quentin Tarantino', 8.9)
) AS dados(titulo, descricao, ano_lancamento, duracao_minutos, genero, diretor, avaliacao)
WHERE NOT EXISTS (SELECT 1 FROM filmes LIMIT 1);

-- Função para atualizar data_atualizacao automaticamente
CREATE OR REPLACE FUNCTION update_data_atualizacao()
RETURNS TRIGGER AS $$
BEGIN
    NEW.data_atualizacao = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger para atualizar automaticamente data_atualizacao
DROP TRIGGER IF EXISTS trigger_update_data_atualizacao ON filmes;
CREATE TRIGGER trigger_update_data_atualizacao
    BEFORE UPDATE ON filmes
    FOR EACH ROW
    EXECUTE FUNCTION update_data_atualizacao();