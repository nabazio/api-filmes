# Módulo 2: Conexão com Banco de Dados e Estruturas
## 🗄️ Integrando PostgreSQL e Criando Modelos

### 📖 Objetivos do Módulo
- Conectar a API Go com banco PostgreSQL
- Criar estruturas (structs) para representar dados
- Implementar operações de leitura do banco
- Organizar código em pacotes especializados
- Implementar tratamento robusto de erros
- Substituir dados estáticos por consultas reais

---

## 🧠 Conceitos Fundamentais

### O que são Structs em Go?
Structs são tipos personalizados que agrupam dados relacionados. É como uma "forma" ou "molde" para seus dados:

```go
type Pessoa struct {
    Nome  string
    Idade int
    Email string
}
```

### O que são Tags JSON?
Tags JSON dizem ao Go como converter structs para JSON e vice-versa:

```go
type Filme struct {
    ID     int    `json:"id"`           // Campo "id" no JSON
    Titulo string `json:"titulo"`       // Campo "titulo" no JSON
    Ano    int    `json:"ano_lancamento"` // Mapeia para "ano_lancamento"
}
```

### Package Database/SQL
É a biblioteca padrão do Go para trabalhar com bancos SQL. Características:
- **Interface uniforme**: Funciona com PostgreSQL, MySQL, SQLite, etc.
- **Prepared statements**: Previne SQL injection automaticamente
- **Connection pooling**: Gerencia conexões automaticamente
- **Context support**: Permite cancelar operações longas

---

## 📁 Evolução da Estrutura do Projeto

### Estrutura Anterior (Módulo 1):
```
api-filmes/
├── cmd/server/main.go
└── go.mod
```

### Nova Estrutura (Módulo 2):
```
api-filmes/
├── cmd/
│   └── server/
│       └── main.go          # Servidor principal
├── internal/
│   ├── models/
│   │   └── filme.go         # ✨ Estruturas de dados
│   ├── database/
│   │   └── conexao.go       # ✨ Conexão e operações de banco
│   └── config/
│       └── config.go        # ✨ Configurações da aplicação
└── go.mod
```

### Por que essa organização?

**🔒 internal/**: Código que não pode ser importado por outros projetos
- **Segurança**: Evita que código interno seja usado indevidamente
- **Encapsulamento**: Mantém detalhes de implementação privados

**📊 models/**: Estruturas que representam dados
- **Centralização**: Todos os tipos de dados em um lugar
- **Reutilização**: Mesmas structs usadas em toda aplicação
- **Documentação**: Serve como documentação dos dados

**🗄️ database/**: Operações de banco de dados
- **Abstração**: Esconde detalhes de SQL do resto da aplicação
- **Testabilidade**: Facilita criação de mocks para testes
- **Manutenibilidade**: Mudanças no banco ficam isoladas

**⚙️ config/**: Configurações da aplicação
- **Flexibilidade**: Fácil mudança entre ambientes (dev, prod)
- **Segurança**: Suporte a variáveis de ambiente
- **Centralização**: Todas as configurações em um lugar

---

## ⚙️ Sistema de Configuração

### Arquivo: `internal/config/config.go`

```go
package config

import (
    "fmt"
    "os"
)

// ConfiguracaoBanco contém as informações de conexão com o banco
type ConfiguracaoBanco struct {
    Host     string
    Porta    string
    Usuario  string
    Senha    string
    NomeBanco string
    SSLMode  string
}
```

**Explicação da Struct:**
- **Host**: Endereço do servidor PostgreSQL (ex: localhost)
- **Porta**: Porta de conexão (padrão PostgreSQL: 5432)
- **Usuario/Senha**: Credenciais de acesso
- **NomeBanco**: Nome do banco de dados
- **SSLMode**: Nível de segurança SSL (disable, require, etc.)

### Função de Configuração

```go
func ObterConfiguracaoBanco() *ConfiguracaoBanco {
    return &ConfiguracaoBanco{
        Host:      obterVariavelOuPadrao("DB_HOST", "localhost"),
        Porta:     obterVariavelOuPadrao("DB_PORT", "5432"),
        Usuario:   obterVariavelOuPadrao("DB_USER", "postgres"),
        Senha:     obterVariavelOuPadrao("DB_PASSWORD", "postgres"),
        NomeBanco: obterVariavelOuPadrao("DB_NAME", "api_filmes"),
        SSLMode:   obterVariavelOuPadrao("DB_SSLMODE", "disable"),
    }
}
```

**Como funciona:**
1. Primeiro tenta ler variável de ambiente (ex: `DB_HOST`)
2. Se não existir, usa valor padrão (ex: "localhost")
3. Permite configuração flexível sem alterar código

### String de Conexão

```go
func (c *ConfiguracaoBanco) StringConexao() string {
    return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Host, c.Porta, c.Usuario, c.Senha, c.NomeBanco, c.SSLMode)
}
```

**Resultado típico:**
```
host=localhost port=5432 user=postgres password=postgres dbname=api_filmes sslmode=disable
```

### Função Helper

```go
func obterVariavelOuPadrao(chave, valorPadrao string) string {
    if valor := os.Getenv(chave); valor != "" {
        return valor
    }
    return valorPadrao
}
```

**Como usar variáveis de ambiente:**
```bash
# No terminal (Linux/Mac)
export DB_HOST="192.168.1.100"
export DB_PASSWORD="senha_segura"

# No Windows
set DB_HOST=192.168.1.100
set DB_PASSWORD=senha_segura

# Executar aplicação
go run cmd/server/main.go
```

---

## 📊 Modelos de Dados

### Arquivo: `internal/models/filme.go`

### Struct Principal: Filme

```go
type Filme struct {
    ID               int       `json:"id"`
    Titulo           string    `json:"titulo"`
    Descricao        string    `json:"descricao"`
    AnoLancamento    int       `json:"ano_lancamento"`
    DuracaoMinutos   int       `json:"duracao_minutos"`
    Genero           string    `json:"genero"`
    Diretor          string    `json:"diretor"`
    Avaliacao        float64   `json:"avaliacao"`
    DataCriacao      time.Time `json:"data_criacao"`
    DataAtualizacao  time.Time `json:"data_atualizacao"`
}
```

**Mapeamento Campo → JSON:**
- `ID` → `"id"`
- `AnoLancamento` → `"ano_lancamento"`
- `DataCriacao` → `"data_criacao"`

**Tipos de Dados Explicados:**
- `int`: Números inteiros (1, 2, 3...)
- `string`: Texto ("O Poderoso Chefão")
- `float64`: Números decimais (9.2, 8.5)
- `time.Time`: Data e hora

### Struct para Listagens: FilmeResumo

```go
type FilmeResumo struct {
    ID            int     `json:"id"`
    Titulo        string  `json:"titulo"`
    AnoLancamento int     `json:"ano_lancamento"`
    Genero        string  `json:"genero"`
    Diretor       string  `json:"diretor"`
    Avaliacao     float64 `json:"avaliacao"`
}
```

**Por que uma struct separada?**
- **Performance**: Menos dados trafegados na rede
- **UI/UX**: Listagens não precisam de todos os detalhes
- **Flexibilidade**: Pode ter campos diferentes da struct principal

### Structs de Resposta

```go
// Para respostas de listagem
type RespostaFilmes struct {
    Filmes []FilmeResumo `json:"filmes"`
    Total  int           `json:"total"`
}

// Para respostas de erro
type RespostaErro struct {
    Erro    string `json:"erro"`
    Codigo  int    `json:"codigo"`
    Detalhes string `json:"detalhes,omitempty"`
}
```

**Tag `omitempty`:**
- Omite o campo do JSON se estiver vazio
- `Detalhes` só aparece se tiver valor

**Exemplo de uso:**
```json
{
    "filmes": [...],
    "total": 3
}

{
    "erro": "Filme não encontrado",
    "codigo": 404
}
```

---

## 🗄️ Camada de Banco de Dados

### Arquivo: `internal/database/conexao.go`

### Struct do Banco

```go
type BancoDados struct {
    conexao *sql.DB
}
```

**Por que encapsular `*sql.DB`?**
- **Abstração**: Esconde detalhes da biblioteca SQL
- **Extensibilidade**: Fácil adicionar cache, logs, métricas
- **Testabilidade**: Pode criar implementações mock

### Estabelecendo Conexão

```go
func NovaConexao() (*BancoDados, error) {
    configuracao := config.ObterConfiguracaoBanco()
    
    fmt.Println("🔌 Conectando ao banco de dados...")
    fmt.Printf("📍 Host: %s:%s | Banco: %s\n", 
        configuracao.Host, configuracao.Porta, configuracao.NomeBanco)
    
    conexao, err := sql.Open("postgres", configuracao.StringConexao())
    if err != nil {
        return nil, fmt.Errorf("erro ao abrir conexão: %v", err)
    }
    
    // Testar a conexão
    if err := conexao.Ping(); err != nil {
        return nil, fmt.Errorf("erro ao conectar com banco: %v", err)
    }
    
    fmt.Println("✅ Conexão com banco estabelecida com sucesso!")
    
    return &BancoDados{conexao: conexao}, nil
}
```

**Passos da conexão:**
1. **sql.Open()**: Cria pool de conexões (não conecta ainda)
2. **Ping()**: Testa se realmente consegue conectar
3. **Error handling**: Retorna erro detalhado se falhar

**Por que usar `fmt.Errorf()`?**
- Adiciona contexto ao erro original
- Facilita debugging
- Mantém stack trace

### Fechando Conexão

```go
func (bd *BancoDados) Fechar() error {
    if bd.conexao != nil {
        return bd.conexao.Close()
    }
    return nil
}
```

**Boa prática**: Sempre fechar conexões para evitar memory leaks.

### Operação: Buscar Todos os Filmes

```go
func (bd *BancoDados) BuscarTodosFilmes() ([]models.FilmeResumo, error) {
    query := `
        SELECT id, titulo, ano_lancamento, genero, diretor, avaliacao 
        FROM filmes 
        ORDER BY titulo ASC
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
```

**Explicação detalhada:**

1. **Query SQL**: Busca apenas campos necessários, ordenados por título
2. **bd.conexao.Query()**: Executa query que retorna múltiplas linhas
3. **defer linhas.Close()**: Garante que recursos sejam liberados
4. **linhas.Next()**: Itera sobre cada linha resultado
5. **linhas.Scan()**: Mapeia colunas SQL para campos da struct
6. **append()**: Adiciona filme à lista
7. **linhas.Err()**: Verifica erros que podem ter ocorrido durante iteração

**Por que usar `&filme.ID`?**
- `Scan()` precisa de ponteiros para modificar os valores
- `&` obtém o endereço de memória da variável

### Operação: Buscar por ID

```go
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
```

**Diferenças importantes:**

- **QueryRow()**: Para buscar apenas uma linha
- **$1**: Placeholder parametrizado (previne SQL injection)
- **sql.ErrNoRows**: Erro específico quando não encontra registros
- Retorna **ponteiro** `*models.Filme` (mais eficiente para structs grandes)

**Por que `$1` em vez de concatenar string?**
```go
// ❌ PERIGOSO - Vulnerável a SQL injection
query := "SELECT * FROM filmes WHERE id = " + id

// ✅ SEGURO - Parametrizado
query := "SELECT * FROM filmes WHERE id = $1"
bd.conexao.QueryRow(query, id)
```

### Operação: Contar Filmes

```go
func (bd *BancoDados) ContarFilmes() (int, error) {
    var total int
    
    query := "SELECT COUNT(*) FROM filmes"
    err := bd.conexao.QueryRow(query).Scan(&total)
    
    if err != nil {
        return 0, fmt.Errorf("erro ao contar filmes: %v", err)
    }
    
    return total, nil
}
```

**Uso típico**: Para paginação e informações estatísticas.

---

## 🌐 Servidor HTTP Aprimorado

### Arquivo: `cmd/server/main.go` Atualizado

### Variável Global e Inicialização

```go
// Variável global para o banco (vamos melhorar isso nos próximos módulos)
var bancoDados *database.BancoDados

func main() {
    fmt.Println("🎬 Servidor da API de Filmes iniciando...")
    
    // Conectar ao banco
    var err error
    bancoDados, err = database.NovaConexao()
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
}
```

**Por que variável global?**
- **Simplicidade**: Para aprendizado inicial
- **Acessibilidade**: Todos os handlers podem usar
- **Limitação**: Dificulta testes (melhoraremos nos próximos módulos)

**Defer com função anônima:**
```go
defer func() {
    // Código executado quando main() termina
}()
```

### Novas Rotas

```go
// Configurar rotas
http.HandleFunc("/filmes", manipularFilmes)
http.HandleFunc("/filmes/", manipularFilmeIndividual)  // ✨ NOVO
http.HandleFunc("/", paginaInicial)
```

**Diferença entre `/filmes` e `/filmes/`:**
- `/filmes`: Exatamente essa URL
- `/filmes/`: Essa URL e tudo que começa com ela (ex: `/filmes/1`, `/filmes/abc`)

### Manipulador de Filme Individual

```go
func manipularFilmeIndividual(w http.ResponseWriter, r *http.Request) {
    configurarCabecalhos(w)
    
    // Extrair ID da URL
    caminho := strings.TrimPrefix(r.URL.Path, "/filmes/")
    if caminho == "" {
        enviarErro(w, "ID do filme é obrigatório", http.StatusBadRequest)
        return
    }
    
    id, err := strconv.Atoi(caminho)
    if err != nil {
        enviarErro(w, "ID inválido", http.StatusBadRequest)
        return
    }
    
    switch r.Method {
    case "GET":
        buscarFilmePorID(w, r, id)
    default:
        enviarErro(w, "Método não permitido", http.StatusMethodNotAllowed)
    }
}
```

**Extração de ID da URL:**
1. **strings.TrimPrefix()**: Remove "/filmes/" do início
2. **strconv.Atoi()**: Converte string para inteiro
3. **Validação**: Verifica se conversão deu certo

**Exemplos:**
- URL: `/filmes/123` → caminho: `"123"` → id: `123`
- URL: `/filmes/abc` → caminho: `"abc"` → erro na conversão

### Função de Listagem Real

```go
func listarFilmes(w http.ResponseWriter, r *http.Request) {
    fmt.Println("📋 Buscando lista de filmes...")
    
    filmes, err := bancoDados.BuscarTodosFilmes()
    if err != nil {
        fmt.Printf("❌ Erro ao buscar filmes: %v\n", err)
        enviarErro(w, "Erro interno do servidor", http.StatusInternalServerError)
        return
    }
    
    total, err := bancoDados.ContarFilmes()
    if err != nil {
        fmt.Printf("⚠️ Erro ao contar filmes: %v\n", err)
        total = len(filmes) // Usar tamanho da lista como fallback
    }
    
    resposta := models.RespostaFilmes{
        Filmes: filmes,
        Total:  total,
    }
    
    fmt.Printf("✅ Encontrados %d filmes\n", len(filmes))
    enviarJSON(w, resposta, http.StatusOK)
}
```

**Fluxo de execução:**
1. Log do início da operação
2. Busca filmes no banco
3. Se erro, retorna erro 500
4. Conta total de filmes
5. Se erro no count, usa tamanho da lista
6. Monta resposta estruturada
7. Retorna JSON com dados

### Função de Busca por ID

```go
func buscarFilmePorID(w http.ResponseWriter, r *http.Request, id int) {
    fmt.Printf("🔍 Buscando filme com ID: %d\n", id)
    
    filme, err := bancoDados.BuscarFilmePorID(id)
    if err != nil {
        if strings.Contains(err.Error(), "não encontrado") {
            enviarErro(w, fmt.Sprintf("Filme com ID %d não encontrado", id), http.StatusNotFound)
        } else {
            fmt.Printf("❌ Erro ao buscar filme: %v\n", err)
            enviarErro(w, "Erro interno do servidor", http.StatusInternalServerError)
        }
        return
    }
    
    fmt.Printf("✅ Filme encontrado: %s\n", filme.Titulo)
    enviarJSON(w, filme, http.StatusOK)
}
```

**Tratamento inteligente de erros:**
- **Não encontrado**: HTTP 404 com mensagem específica
- **Erro de banco**: HTTP 500 com mensagem genérica
- **Log detalhado**: Para debugging do desenvolvedor

### Funções Utilitárias

```go
// Configurar cabeçalhos HTTP
func configurarCabecalhos(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
```

**Headers CORS explicados:**
- **Access-Control-Allow-Origin**: Permite requests de qualquer origem
- **Access-Control-Allow-Methods**: Métodos HTTP permitidos
- **Access-Control-Allow-Headers**: Headers que frontend pode enviar

```go
// Enviar resposta JSON
func enviarJSON(w http.ResponseWriter, dados interface{}, status int) {
    w.WriteHeader(status)
    
    if err := json.NewEncoder(w).Encode(dados); err != nil {
        fmt.Printf("❌ Erro ao codificar JSON: %v\n", err)
        http.Error(w, "Erro interno", http.StatusInternalServerError)
    }
}
```

**json.NewEncoder(w).Encode():**
- Converte struct para JSON automaticamente
- Escreve diretamente na resposta HTTP
- Mais eficiente que `json.Marshal()`

```go
// Enviar erro em formato JSON
func enviarErro(w http.ResponseWriter, mensagem string, status int) {
    erro := models.RespostaErro{
        Erro:   mensagem,
        Codigo: status,
    }
    
    enviarJSON(w, erro, status)
}
```

**Padronização de erros:**
- Sempre retorna JSON (nunca texto simples)
- Inclui código HTTP no corpo da resposta
- Facilita tratamento no frontend

---

## 🧪 Testes Completos

### 1. Preparação do Ambiente

```bash
# Verificar se banco está rodando
# No DBeaver, executar:
SELECT COUNT(*) FROM filmes;

# Deve retornar 3 (filmes inseridos no Módulo 1)
```

### 2. Executar Servidor

```bash
cd api-filmes
go run cmd/server/main.go
```

**Saída esperada:**
```
🎬 Servidor da API de Filmes iniciando...
🔌 Conectando ao banco de dados...
📍 Host: localhost:5432 | Banco: api_filmes
✅ Conexão com banco estabelecida com sucesso!
🚀 Servidor rodando em http://localhost:8080
```

### 3. Testes no Postman

#### Teste 1: Página Inicial
- **Método**: GET
- **URL**: `http://localhost:8080/`
- **Resultado esperado**:
```json
{
    "mensagem": "🎬 Bem-vindo à API de Filmes!",
    "versao": "1.0.0",
    "endpoints": [
        "GET /filmes - Lista todos os filmes",
        "GET /filmes/{id} - Busca filme por ID"
    ]
}
```

#### Teste 2: Lista Todos os Filmes
- **Método**: GET
- **URL**: `http://localhost:8080/filmes`
- **Resultado esperado**:
```json
{
    "filmes": [
        {
            "id": 2,
            "titulo": "Cidade de Deus",
            "ano_lancamento": 2002,
            "genero": "Drama",
            "diretor": "Fernando Meirelles",
            "avaliacao": 8.6
        },
        {
            "id": 1,
            "titulo": "O Poderoso Chefão",
            "ano_lancamento": 1972,
            "genero": "Drama",
            "diretor": "Francis Ford Coppola",
            "avaliacao": 9.2
        },
        {
            "id": 3,
            "titulo": "Vingadores: Ultimato",
            "ano_lancamento": 2019,
            "genero": "Ação",
            "diretor": "Anthony e Joe Russo",
            "avaliacao": 8.4
        }
    ],
    "total": 3
}
```

#### Teste 3: Filme por ID (Existente)
- **Método**: GET
- **URL**: `http://localhost:8080/filmes/1`
- **Resultado esperado**:
```json
{
    "id": 1,
    "titulo": "O Poderoso Chefão",
    "descricao": "A saga de uma família mafiosa italiana nos Estados Unidos",
    "ano_lancamento": 1972,
    "duracao_minutos": 175,
    "genero": "Drama",
    "diretor": "Francis Ford Coppola",
    "avaliacao": 9.2,
    "data_criacao": "2024-01-15T10:30:00Z",
    "data_atualizacao": "2024-01-15T10:30:00Z"
}
```

#### Teste 4: Filme por ID (Inexistente)
- **Método**: GET
- **URL**: `http://localhost:8080/filmes/999`
- **Status**: 404 Not Found
- **Resultado esperado**:
```json
{
    "erro": "Filme com ID 999 não encontrado",
    "codigo": 404
}
```

#### Teste 5: ID Inválido
- **Método**: GET
- **URL**: `http://localhost:8080/filmes/abc`
- **Status**: 400 Bad Request
- **Resultado esperado**:
```json
{
    "erro": "ID inválido",
    "codigo": 400
}
```

#### Teste 6: Método Não Permitido
- **Método**: POST
- **URL**: `http://localhost:8080/filmes/1`
- **Status**: 405 Method Not Allowed
- **Resultado esperado**:
```json
{
    "erro": "Método não permitido",
    "codigo": 405
}
```

---

## 🎓 Conceitos Aprendidos

### 1. Arquitetura em Camadas
- **Config**: Configurações centralizadas
- **Models**: Estruturas de dados
- **Database**: Acesso aos dados
- **Handlers**: Lógica HTTP

### 2. Structs e JSON
- Tags para mapeamento automático
- Diferentes structs para diferentes necessidades
- Conversão bidirecional Go ↔ JSON

### 3. Database/SQL
- Connection pooling automático
- Prepared statements para segurança
- Diferentes métodos: Query vs QueryRow
- Tratamento específico de erros

### 4. HTTP Melhorado
- Roteamento por padrão de URL
- Headers CORS para frontend
- Status codes apropriados
- Respostas JSON estruturadas

### 5. Tratamento de Erros
- Errors com contexto usando `fmt.Errorf()`
- Distinção entre erro de negócio vs sistema
- Logs para developer vs mensagens para usuário

---

## 🔧 Troubleshooting

### Problema: "conexão recusada"
```bash
# Verificar se PostgreSQL está rodando
sudo systemctl status postgresql

# Ou no Windows
net start postgresql

# Verificar conexão no DBeaver primeiro
```

### Problema: "tabela não existe"
```sql
-- Verificar se banco e tabela existem
\c api_filmes
\dt
SELECT * FROM filmes;
```

### Problema: "driver não encontrado"
```bash
# Verificar se dependência foi instalada
go mod tidy

# Se necessário, reinstalar
go get github.com/lib/pq
```

### Problema: JSON malformado
- Verificar tags JSON nas structs
- Usar ferramenta online para validar JSON
- Verificar se todos os campos estão sendo populados

### Problema: "método não permitido"
- Verificar método HTTP no Postman (GET/POST/etc)
- Confirmar URL exata
- Verificar se handler está registrado corretamente

---

## 📊 Métricas de Sucesso

Ao final deste módulo, você deve ser capaz de:

- [ ] Conectar Go com PostgreSQL sem erros
- [ ] Criar structs com tags JSON apropriadas
- [ ] Implementar operações de leitura do banco
- [ ] Organizar código em pacotes lógicos
- [ ] Tratar diferentes tipos de erro adequadamente
- [ ] Testar API com dados reais via Postman
- [ ] Entender fluxo: Request → Handler → Database → Response

---

## 🚀 Comparação: Antes vs Depois

### Módulo 1 (Dados Estáticos):
```go
func listarFilmes(w http.ResponseWriter, r *http.Request) {
    filmesJson := `{"filmes": [...]}`  // String fixa
    w.Write([]byte(filmesJson))
}
```

### Módulo 2 (Dados Dinâmicos):
```go
func listarFilmes(w http.ResponseWriter, r *http.Request) {
    filmes, err := bancoDados.BuscarTodosFilmes()  // Busca real no banco
    if err != nil {
        enviarErro(w, "Erro interno", 500)         // Tratamento de erro
        return
    }
    
    resposta := models.RespostaFilmes{             // Struct tipada
        Filmes: filmes,
        Total:  len(filmes),
    }
    
    enviarJSON(w, resposta, http.StatusOK)         // Conversão automática
}
```

**Evolução alcançada:**
- ✅ **Dados dinâmicos** em vez de strings fixas
- ✅ **Tratamento de erro robusto** com logs e status codes
- ✅ **Tipagem forte** com structs em vez de strings
- ✅ **Conversão automática** JSON com reflection
- ✅ **Separação de responsabilidades** entre camadas

---

## 🏗️ Padrões de Arquitetura Implementados

### 1. Repository Pattern (Rudimentar)
```go
// interface implícita para operações de banco
type FilmeRepository interface {
    BuscarTodosFilmes() ([]FilmeResumo, error)
    BuscarFilmePorID(id int) (*Filme, error)
    ContarFilmes() (int, error)
}

// BancoDados implementa implicitamente FilmeRepository
type BancoDados struct {
    conexao *sql.DB
}
```

### 2. Data Transfer Object (DTO)
```go
// FilmeResumo = DTO para listagens
type FilmeResumo struct {
    ID            int     `json:"id"`
    Titulo        string  `json:"titulo"`
    // ... apenas campos necessários
}

// Filme = Entidade completa para detalhes
type Filme struct {
    // ... todos os campos
}
```

### 3. Response Wrapper Pattern
```go
// Padronização de respostas
type RespostaFilmes struct {
    Filmes []FilmeResumo `json:"filmes"`
    Total  int           `json:"total"`
}

type RespostaErro struct {
    Erro   string `json:"erro"`
    Codigo int    `json:"codigo"`
}
```

---

## 🔍 Análise de Performance

### Queries Otimizadas

**Lista de filmes:**
```sql
-- ✅ Busca apenas campos necessários
SELECT id, titulo, ano_lancamento, genero, diretor, avaliacao 
FROM filmes 
ORDER BY titulo ASC;

-- ❌ Evitamos SELECT *
-- SELECT * FROM filmes;  -- Trafega dados desnecessários
```

**Busca por ID:**
```sql
-- ✅ Usa índice da chave primária
SELECT * FROM filmes WHERE id = $1;

-- ✅ Parametrizada (previne SQL injection)
-- ❌ Evitamos concatenação
-- "SELECT * FROM filmes WHERE id = " + id
```

### Connection Pooling
```go
// Go gerencia automaticamente pool de conexões
db, err := sql.Open("postgres", connectionString)

// Configurações padrão otimizadas:
// - MaxOpenConns: ilimitado
// - MaxIdleConns: 2
// - ConnMaxLifetime: ilimitado
```

### Memory Management
```go
// ✅ defer garante limpeza de recursos
defer linhas.Close()

// ✅ Slices crescem dinamicamente
var filmes []models.FilmeResumo
filmes = append(filmes, filme)

// ✅ Ponteiros para structs grandes
func BuscarFilmePorID(id int) (*models.Filme, error)
```

---

## 🛡️ Aspectos de Segurança

### 1. SQL Injection Prevention
```go
// ✅ SEGURO - Prepared statement
query := "SELECT * FROM filmes WHERE id = $1"
err := db.QueryRow(query, id).Scan(...)

// ❌ VULNERÁVEL - String concatenation
// query := "SELECT * FROM filmes WHERE id = " + userInput
```

### 2. Error Information Disclosure
```go
// ✅ SEGURO - Erro genérico para usuário
if err != nil {
    log.Printf("Erro específico: %v", err)  // Log interno
    enviarErro(w, "Erro interno", 500)      // Resposta genérica
}

// ❌ PERIGOSO - Vazar detalhes do sistema
// w.Write([]byte(err.Error()))  // Pode expor paths, conexões, etc.
```

### 3. Input Validation
```go
// ✅ Validação de tipo
id, err := strconv.Atoi(caminho)
if err != nil {
    enviarErro(w, "ID inválido", 400)
    return
}

// ✅ Validação de range (implemente quando necessário)
if id <= 0 {
    enviarErro(w, "ID deve ser positivo", 400)
    return
}
```

### 4. CORS Configuration
```go
// ✅ Para desenvolvimento - permite todas as origens
w.Header().Set("Access-Control-Allow-Origin", "*")

// 🚨 Para produção - especificar domínios
// w.Header().Set("Access-Control-Allow-Origin", "https://meusite.com")
```

---

## 📈 Monitoramento e Debugging

### Logs Implementados
```go
// Início de operações
fmt.Println("📋 Buscando lista de filmes...")

// Sucessos
fmt.Printf("✅ Encontrados %d filmes\n", len(filmes))

// Erros com contexto
fmt.Printf("❌ Erro ao buscar filmes: %v\n", err)

// Informações de conexão
fmt.Printf("📍 Host: %s:%s | Banco: %s\n", host, porta, banco)
```

### Métricas Básicas
```go
// Total de registros
total, err := bancoDados.ContarFilmes()

// Tempo de resposta (adicione quando necessário)
start := time.Now()
// ... operação ...
duration := time.Since(start)
fmt.Printf("⏱️ Operação levou: %v\n", duration)
```

### Health Check Básico
```go
// Teste de conexão
if err := conexao.Ping(); err != nil {
    return nil, fmt.Errorf("erro ao conectar com banco: %v", err)
}
```

---

## 🔄 Fluxo de Dados Completo

### Request → Response Flow

```
1. Cliente (Postman)
   ↓ GET /filmes/1
   
2. Go HTTP Server
   ↓ manipularFilmeIndividual()
   
3. URL Parsing
   ↓ extrair ID da URL
   
4. Validation
   ↓ validar se ID é inteiro
   
5. Database Layer
   ↓ bancoDados.BuscarFilmePorID(1)
   
6. PostgreSQL
   ↓ SELECT * FROM filmes WHERE id = $1
   
7. Row Scanning
   ↓ mapear colunas → struct Filme
   
8. JSON Encoding
   ↓ converter struct → JSON
   
9. HTTP Response
   ↓ status code + headers + body
   
10. Cliente recebe JSON
```

### Error Flow

```
1. Erro no banco
   ↓ connection refused / query error
   
2. Database Layer
   ↓ return nil, fmt.Errorf("contexto: %v", err)
   
3. Handler Layer  
   ↓ log específico + resposta genérica
   
4. JSON Error Response
   ↓ {"erro": "Erro interno", "codigo": 500}
   
5. Cliente recebe erro estruturado
```

---

## 🎯 Preparação para Módulo 3

### O que já temos funcionando:
- ✅ Conexão estável com PostgreSQL
- ✅ Operações de leitura (GET)
- ✅ Estruturas de dados bem definidas
- ✅ Tratamento básico de erros
- ✅ Respostas JSON padronizadas

### O que precisamos adicionar:
- 🔜 Operações de escrita (POST, PUT, DELETE)
- 🔜 Validação robusta de dados de entrada
- 🔜 Middleware para logs e autenticação
- 🔜 Handlers organizados em arquivos separados
- 🔜 Testes automatizados

### Conceitos que vamos aprender:
- **Request Body Parsing**: Como ler JSON do cliente
- **Data Validation**: Validar dados antes de salvar
- **HTTP Methods**: Implementar POST, PUT, DELETE
- **Middleware Pattern**: Código que roda antes/depois dos handlers
- **Error Handling**: Tratamento mais sofisticado
- **Code Organization**: Separar responsabilidades melhor

---

## 📚 Referências e Estudo Adicional

### Documentação Oficial
- [database/sql Package](https://pkg.go.dev/database/sql) - Documentação da biblioteca SQL
- [encoding/json Package](https://pkg.go.dev/encoding/json) - Conversão JSON
- [net/http Package](https://pkg.go.dev/net/http) - Servidor HTTP

### Tutoriais Recomendados
- [Go by Example - JSON](https://gobyexample.com/json)
- [Go by Example - Structs](https://gobyexample.com/structs)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/current/)

### Ferramentas Úteis
- **DBeaver**: Interface gráfica para banco de dados
- **Postman**: Teste de APIs
- **Go Playground**: Testar código Go online
- **JSON Formatter**: Validar e formatar JSON

### Próximos Estudos
- Padrões de arquitetura em Go
- Testes unitários com testify
- Dockerização de aplicações Go
- Deploy em produção

---

## ✅ Checklist Final do Módulo 2

Antes de prosseguir para o Módulo 3, verifique:

### Configuração:
- [ ] Banco `api_filmes` funcionando no PostgreSQL
- [ ] Tabela `filmes` com dados de exemplo
- [ ] Dependência `github.com/lib/pq` instalada
- [ ] Estrutura de pastas criada corretamente

### Código:
- [ ] Arquivo `config/config.go` criado e funcionando
- [ ] Arquivo `models/filme.go` com todas as structs
- [ ] Arquivo `database/conexao.go` com operações de banco
- [ ] Arquivo `main.go` atualizado com nova lógica

### Testes:
- [ ] Servidor inicia sem erros de conexão
- [ ] GET `/` retorna informações da API
- [ ] GET `/filmes` lista todos os filmes do banco
- [ ] GET `/filmes/1` retorna filme específico
- [ ] GET `/filmes/999` retorna erro 404
- [ ] GET `/filmes/abc` retorna erro 400
- [ ] Logs aparecem no console durante operações

### Compreensão:
- [ ] Entendo diferença entre Query() e QueryRow()
- [ ] Sei como structs se convertem para JSON
- [ ] Compreendo fluxo request → database → response
- [ ] Reconheço importância de prepared statements
- [ ] Entendo papel de cada camada (config, models, database, handlers)

---

**🚀 Parabéns! Você completou o Módulo 2 e agora tem uma API funcional conectada ao banco de dados PostgreSQL!**

**No Módulo 3, vamos implementar operações de criação, atualização e exclusão de filmes, além de melhorar a organização do código e adicionar validações robustas.**