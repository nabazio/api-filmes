# Módulo 1: Fundações da API em Go
## 🎯 Configuração Inicial e Estrutura Base

### 📖 Objetivos do Módulo
- Criar a estrutura base do projeto Go
- Configurar um servidor HTTP básico
- Implementar rotas simples
- Entender conceitos fundamentais de APIs REST
- Preparar o ambiente de desenvolvimento

---

## 🧠 Conceitos Fundamentais

### O que é uma API REST?
Uma API REST (Representational State Transfer) é um conjunto de regras e convenções para criar serviços web. Usando a analogia do restaurante:

- **Cliente (Frontend/Postman)** → Faz pedidos
- **Garçom (API)** → Intermedia e organiza os pedidos
- **Cozinha (Servidor/Banco)** → Processa e prepara a resposta

### Estrutura de uma Requisição HTTP
```
Método + URL + Cabeçalhos + Corpo (opcional)
GET /filmes HTTP/1.1
Host: localhost:8080
Content-Type: application/json
```

### Principais Métodos HTTP
- **GET** → Buscar dados (como "ver o cardápio")
- **POST** → Criar novos dados (como "fazer um pedido")
- **PUT** → Atualizar dados completos
- **DELETE** → Remover dados

---

## 📁 Estrutura do Projeto

### Por que essa organização?
```
api-filmes/
├── cmd/              # Comandos executáveis (pontos de entrada)
│   └── server/       # Servidor principal
├── internal/         # Código interno (não pode ser importado por outros projetos)
│   ├── handlers/     # Funções que processam requisições HTTP
│   ├── models/       # Estruturas de dados
│   ├── database/     # Conexão e operações de banco
│   └── config/       # Configurações da aplicação
├── pkg/              # Código que pode ser reutilizado por outros projetos
└── docs/             # Documentação
```

**Vantagens dessa estrutura:**
- **Escalabilidade**: Fácil de adicionar novos recursos
- **Manutenibilidade**: Cada coisa tem seu lugar
- **Testabilidade**: Código bem separado é mais fácil de testar
- **Padrão Go**: Segue as convenções da linguagem

---

## 🗄️ Preparação do Banco de Dados

### 1. Criando o Banco no DBeaver

**Passo a passo:**
1. Abra o DBeaver
2. Conecte-se ao seu servidor PostgreSQL
3. Execute os comandos SQL:

```sql
-- Criação do banco
CREATE DATABASE api_filmes;
```

### 2. Estrutura da Tabela

```sql
-- Conectar ao banco criado
\c api_filmes;

-- Criar tabela de filmes
CREATE TABLE filmes (
    id SERIAL PRIMARY KEY,                    -- ID auto-incremento
    titulo VARCHAR(255) NOT NULL,             -- Nome do filme (obrigatório)
    descricao TEXT,                          -- Sinopse (opcional)
    ano_lancamento INTEGER NOT NULL,          -- Ano (obrigatório)
    duracao_minutos INTEGER,                 -- Duração em minutos
    genero VARCHAR(100),                     -- Categoria do filme
    diretor VARCHAR(255),                    -- Nome do diretor
    avaliacao DECIMAL(3,1) CHECK (avaliacao >= 0 AND avaliacao <= 10), -- Nota de 0 a 10
    data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    -- Quando foi criado
    data_atualizacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Última atualização
);
```

**Explicando os tipos de dados:**
- `SERIAL`: Número inteiro que incrementa automaticamente
- `VARCHAR(255)`: Texto com limite de 255 caracteres
- `TEXT`: Texto sem limite específico
- `INTEGER`: Número inteiro
- `DECIMAL(3,1)`: Número decimal com 3 dígitos total, 1 após a vírgula
- `TIMESTAMP`: Data e hora
- `CHECK`: Validação para garantir que avaliação está entre 0 e 10

### 3. Dados de Exemplo

```sql
INSERT INTO filmes (titulo, descricao, ano_lancamento, duracao_minutos, genero, diretor, avaliacao) VALUES
('O Poderoso Chefão', 'A saga de uma família mafiosa italiana nos Estados Unidos', 1972, 175, 'Drama', 'Francis Ford Coppola', 9.2),
('Cidade de Deus', 'Retrato da violência urbana no Rio de Janeiro', 2002, 130, 'Drama', 'Fernando Meirelles', 8.6),
('Vingadores: Ultimato', 'Os heróis se unem para derrotar Thanos', 2019, 181, 'Ação', 'Anthony e Joe Russo', 8.4);
```

---

## 🚀 Configuração do Projeto Go

### 1. Inicializando o Módulo

```bash
# Criar pasta do projeto
mkdir api-filmes
cd api-filmes

# Inicializar módulo Go
go mod init api-filmes
```

**O que faz `go mod init`?**
- Cria o arquivo `go.mod` que gerencia dependências
- Define o nome do módulo (usado para importações)
- É como o `package.json` do Node.js

### 2. Instalando Dependências

```bash
# Driver para PostgreSQL
go get github.com/lib/pq
```

**O que é um driver de banco?**
É um código que permite ao Go "conversar" com o PostgreSQL. Cada banco tem seu próprio driver.

---

## 💻 Código Explicado

### Arquivo main.go Completo

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    fmt.Println("🎬 Servidor da API de Filmes iniciando...")
    
    // Configurar rota básica
    http.HandleFunc("/filmes", manipularFilmes)
    http.HandleFunc("/", paginaInicial)
    
    // Iniciar servidor
    porta := ":8080"
    fmt.Printf("🚀 Servidor rodando em http://localhost%s\n", porta)
    
    if err := http.ListenAndServe(porta, nil); err != nil {
        log.Fatal("❌ Erro ao iniciar servidor:", err)
    }
}
```

**Explicação linha por linha:**

1. `package main`: Define que este é o pacote principal (executável)
2. `import`: Importa bibliotecas necessárias
3. `http.HandleFunc()`: Registra uma função para uma rota específica
4. `http.ListenAndServe()`: Inicia o servidor na porta especificada

### Função de Página Inicial

```go
func paginaInicial(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    resposta := `{
        "mensagem": "🎬 Bem-vindo à API de Filmes!",
        "endpoints": [
            "GET /filmes - Lista todos os filmes"
        ]
    }`
    w.Write([]byte(resposta))
}
```

**Parâmetros da função:**
- `w http.ResponseWriter`: Para escrever a resposta de volta ao cliente
- `r *http.Request`: Contém informações sobre a requisição recebida

**O que cada linha faz:**
1. Define o tipo de conteúdo como JSON
2. Cria uma string JSON com informações da API
3. Envia a resposta para o cliente

### Manipulador de Rotas

```go
func manipularFilmes(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    
    switch r.Method {
    case "GET":
        listarFilmes(w, r)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
        w.Write([]byte(`{"erro": "Método não permitido"}`))
    }
}
```

**Por que usar switch para métodos HTTP?**
- Diferentes métodos fazem ações diferentes
- GET para buscar, POST para criar, etc.
- É uma forma organizada de separar a lógica

### Listagem Temporária de Filmes

```go
func listarFilmes(w http.ResponseWriter, r *http.Request) {
    filmesJson := `{
        "filmes": [
            {
                "id": 1,
                "titulo": "O Poderoso Chefão",
                "ano": 1972,
                "genero": "Drama"
            },
            {
                "id": 2,
                "titulo": "Cidade de Deus", 
                "ano": 2002,
                "genero": "Drama"
            }
        ]
    }`
    
    w.Write([]byte(filmesJson))
}
```

**Por que dados temporários?**
- Primeiro testamos se a estrutura funciona
- Depois conectamos com o banco real
- É uma abordagem incremental e segura

---

## 🧪 Testando a API

### 1. Executando o Servidor

```bash
# No terminal, dentro da pasta api-filmes
go run cmd/server/main.go
```

**O que você deve ver:**
```
🎬 Servidor da API de Filmes iniciando...
🚀 Servidor rodando em http://localhost:8080
```

### 2. Testes no Postman

#### Teste 1: Página Inicial
- **Método**: GET
- **URL**: `http://localhost:8080/`
- **Resultado esperado**:
```json
{
    "mensagem": "🎬 Bem-vindo à API de Filmes!",
    "endpoints": [
        "GET /filmes - Lista todos os filmes"
    ]
}
```

#### Teste 2: Lista de Filmes
- **Método**: GET
- **URL**: `http://localhost:8080/filmes`
- **Resultado esperado**:
```json
{
    "filmes": [
        {
            "id": 1,
            "titulo": "O Poderoso Chefão",
            "ano": 1972,
            "genero": "Drama"
        },
        {
            "id": 2,
            "titulo": "Cidade de Deus",
            "ano": 2002,
            "genero": "Drama"
        }
    ]
}
```

#### Teste 3: Método Não Permitido
- **Método**: POST
- **URL**: `http://localhost:8080/filmes`
- **Resultado esperado**:
```json
{
    "erro": "Método não permitido"
}
```

---

## 🎓 Conceitos Aprendidos

### 1. Servidor HTTP em Go
- Go tem uma biblioteca HTTP robusta e simples
- `http.HandleFunc()` conecta URLs a funções
- `http.ListenAndServe()` mantém o servidor rodando

### 2. Manipuladores de Requisição
- Toda função de manipulador recebe `(w, r)`
- `w` é para escrever respostas
- `r` contém dados da requisição

### 3. Roteamento Básico
- Cada rota aponta para uma função específica
- Podemos verificar o método HTTP com `r.Method`
- Switch é uma forma elegante de lidar com diferentes métodos

### 4. Respostas JSON
- APIs modernas usam JSON para trocar dados
- Header `Content-Type` informa o tipo de conteúdo
- `w.Write()` envia dados para o cliente

### 5. Códigos de Status HTTP
- `200 OK`: Tudo funcionou perfeitamente
- `405 Method Not Allowed`: Método HTTP não suportado
- `500 Internal Server Error`: Erro no servidor

---

## 🔧 Troubleshooting (Resolução de Problemas)

### Problema: "Porta já está em uso"
```bash
# Encontrar processo usando a porta 8080
lsof -i :8080

# Matar o processo (substitua PID pelo número encontrado)
kill -9 PID
```

### Problema: "go: comando não encontrado"
- Verifique se Go está instalado: `go version`
- Se não estiver, baixe em: https://golang.org/dl/

### Problema: Postman não conecta
- Verifique se o servidor está rodando
- Confirme a URL: `http://localhost:8080`
- Verifique se não há firewall bloqueando

### Problema: Erro ao criar banco
- Verifique se PostgreSQL está rodando
- Confirme usuário e senha no DBeaver
- Teste a conexão no DBeaver antes de executar comandos

---

## 📋 Checklist de Conclusão

Antes de ir para o Módulo 2, verifique:

- [ ] Projeto Go inicializado com `go mod init`
- [ ] Banco `api_filmes` criado no PostgreSQL
- [ ] Tabela `filmes` criada com dados de exemplo
- [ ] Servidor roda sem erros em `localhost:8080`
- [ ] Teste GET `/` retorna página inicial
- [ ] Teste GET `/filmes` retorna lista de filmes
- [ ] Teste POST `/filmes` retorna erro 405
- [ ] Entendi como funciona roteamento em Go
- [ ] Sei a diferença entre `w` e `r` nos manipuladores

---

## 🎯 Próximos Passos

No **Módulo 2** você aprenderá:

- ✨ Conectar Go com PostgreSQL
- 🏗️ Criar estruturas (structs) para representar dados
- 🔄 Substituir dados fixos por consultas reais ao banco
- 📊 Implementar busca por ID específico
- 🛡️ Melhorar tratamento de erros

**Dica de Estudo**: Releia este módulo antes de continuar e pratique modificando o código para se familiarizar com a sintaxe do Go!

---

## 📚 Referências Adicionais

- [Tour of Go](https://tour.golang.org/) - Tutorial oficial do Go
- [Effective Go](https://golang.org/doc/effective_go.html) - Boas práticas
- [HTTP Package Documentation](https://pkg.go.dev/net/http) - Documentação do pacote HTTP
- [PostgreSQL Documentation](https://www.postgresql.org/docs/) - Documentação do PostgreSQL