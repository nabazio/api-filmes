# Módulo 3-A: Implementando Operações de Criação (POST)
## 📝 Criando Novos Filmes na API

### 📖 Objetivos do Módulo
- Implementar endpoint POST para criar novos filmes
- Adicionar validação robusta de dados de entrada
- Reorganizar código com handlers especializados
- Implementar middleware para logs e tratamento de erros
- Processar corretamente request body em formato JSON
- Melhorar a arquitetura da aplicação

---

## 🧠 Conceitos Fundamentais

### O que são Métodos HTTP?
Os métodos HTTP definem a **intenção** de uma requisição:

- **GET**: "Quero buscar dados" (não modifica nada)
- **POST**: "Quero criar algo novo"
- **PUT**: "Quero atualizar completamente"
- **DELETE**: "Quero remover"
- **OPTIONS**: "Quero saber que operações posso fazer"

### Request Body vs URL Parameters
```
GET /filmes/123        ← ID na URL (parâmetro)
POST /filmes           ← Dados no body JSON
{
  "titulo": "Novo Filme",
  "ano": 2024
}
```

### O que é Middleware?
Middleware é código que executa **entre** o recebimento da requisição e o processamento final:

```
Request → Middleware 1 → Middleware 2 → Handler → Response
          (Logs)        (CORS)         (Lógica)
```

### Estruturas vs Handlers
- **Struct**: Agrupa dados relacionados
- **Methods em Struct**: Funções que "pertencem" à struct
- **Handler**: Função que processa requisições HTTP

---

## 🏗️ Evolução da Arquitetura

### Estrutura Anterior (Módulo 2):
```
api-filmes/
├── cmd/server/main.go     # Tudo misturado
├── internal/
│   ├── models/
│   ├── database/
│   └── config/
```

### Nova Estrutura (Módulo 3-A):
```
api-filmes/
├── cmd/
│   └── server/
│       └── main.go        # Apenas inicialização
├── internal/
│   ├── handlers/          # ✨ Lógica HTTP separada
│   │   ├── filme_handlers.go
│   │   └── middleware.go
│   ├── models/            # Estruturas de dados
│   ├── database/          # Operações de banco
│   ├── config/            # Configurações
│   └── validators/        # ✨ Validação de dados
│       └── filme_validator.go
```

### Benefícios da Nova Estrutura:

**🎯 Separação de Responsabilidades:**
- `handlers/`: Apenas lógica HTTP
- `validators/`: Apenas validação de dados
- `database/`: Apenas operações de banco
- `models/`: Apenas estruturas de dados

**📈 Escalabilidade:**
- Fácil adicionar novos recursos
- Código organizado por função
- Reutilização de componentes

**🧪 Testabilidade:**
- Cada camada pode ser testada isoladamente
- Mocks mais fáceis de criar
- Testes unitários específicos

---

## 🏭 Sistema de Handlers Organizado

### Arquivo: `internal/handlers/filme_handlers.go`

### Struct FilmeHandler

```go
type FilmeHandler struct {
    bancoDados *database.BancoDados
}

func NovoFilmeHandler(bd *database.BancoDados) *FilmeHandler {
    return &FilmeHandler{
        bancoDados: bd,
    }
}
```

**Por que usar Struct para Handlers?**

**✅ Vantagens:**
- **Encapsulamento**: Todos os métodos relacionados ficam juntos
- **Estado compartilhado**: Conexão de banco disponível para todos os métodos
- **Organização**: Agrupa funcionalidades relacionadas
- **Extensibilidade**: Fácil adicionar novos campos (cache, logger, etc.)

**Comparação com funções globais:**
```go
// ❌ Antes: função global com variável global
var bancoDados *database.BancoDados

func listarFilmes(w http.ResponseWriter, r *http.Request) {
    // usa variável global
}

// ✅ Agora: método em struct
func (fh *FilmeHandler) listarFilmes(w http.ResponseWriter, r *http.Request) {
    // usa fh.bancoDados
}
```

### Método ManipularFilmes

```go
func (fh *FilmeHandler) ManipularFilmes(w http.ResponseWriter, r *http.Request) {
    configurarCabecalhos(w)
    
    switch r.Method {
    case "GET":
        fh.listarFilmes(w, r)
    case "POST":
        fh.criarFilme(w, r)          // ✨ NOVO
    case "OPTIONS":
        w.WriteHeader(http.StatusOK) // ✨ CORS Support
    default:
        enviarErro(w, "Método não permitido", http.StatusMethodNotAllowed)
    }
}
```

**Explicação:**
1. **configurarCabecalhos()**: Define headers padrão (JSON, CORS)
2. **switch r.Method**: Roteia baseado no método HTTP
3. **OPTIONS**: Resposta para preflight requests do CORS
4. **default**: Qualquer método não suportado retorna 405

### Método criarFilme (Coração do Módulo)

```go
func (fh *FilmeHandler) criarFilme(w http.ResponseWriter, r *http.Request) {
    fmt.Println("➕ Criando novo filme...")
    
    // 1. Verificar Content-Type
    if r.Header.Get("Content-Type") != "application/json" {
        enviarErro(w, "Content-Type deve ser application/json", http.StatusBadRequest)
        return
    }
```

**Por que verificar Content-Type?**
- **Segurança**: Evita processamento incorreto de dados
- **Clareza**: Cliente sabe exatamente que formato enviar
- **Robustez**: Evita erros de parsing

```go
    // 2. Decodificar JSON do body
    var novoFilme models.CriarFilme
    decoder := json.NewDecoder(r.Body)
    decoder.DisallowUnknownFields() // Rejeita campos não reconhecidos
    
    if err := decoder.Decode(&novoFilme); err != nil {
        fmt.Printf("❌ Erro ao decodificar JSON: %v\n", err)
        enviarErro(w, "JSON inválido", http.StatusBadRequest)
        return
    }
```

**json.NewDecoder vs json.Unmarshal:**
```go
// ✅ NewDecoder - para http.Request.Body
decoder := json.NewDecoder(r.Body)
decoder.Decode(&struct)

// ✅ Unmarshal - para []byte existente
var data []byte
json.Unmarshal(data, &struct)
```

**DisallowUnknownFields():**
```json
// ❌ Será rejeitado
{
  "titulo": "Filme",
  "campo_inexistente": "valor"
}

// ✅ Será aceito
{
  "titulo": "Filme",
  "ano_lancamento": 2024
}
```

```go
    // 3. Validar dados
    if erros := validators.ValidarCriarFilme(&novoFilme); len(erros) > 0 {
        fmt.Printf("❌ Dados inválidos: %v\n", erros)
        resposta := models.RespostaErro{
            Erro:     "Dados inválidos",
            Codigo:   http.StatusBadRequest,
            Detalhes: strings.Join(erros, "; "),
        }
        enviarJSON(w, resposta, http.StatusBadRequest)
        return
    }
```

**Sistema de Validação:**
- Retorna **slice de erros** (múltiplos problemas)
- **strings.Join()** combina erros em uma string
- **Campo Detalhes** fornece informações específicas

```go
    // 4. Salvar no banco
    filmeID, err := fh.bancoDados.CriarFilme(&novoFilme)
    if err != nil {
        fmt.Printf("❌ Erro ao salvar filme: %v\n", err)
        
        // Verificar se é erro de duplicação
        if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
            enviarErro(w, "Filme com este título já existe", http.StatusConflict)
        } else {
            enviarErro(w, "Erro interno do servidor", http.StatusInternalServerError)
        }
        return
    }
```

**Tratamento Inteligente de Erros:**
- **409 Conflict**: Para violações de unicidade
- **500 Internal Server Error**: Para outros erros de banco
- **Log detalhado**: Para debugging
- **Mensagem amigável**: Para o usuário

```go
    // 5. Buscar filme criado para retornar completo
    filmeCriado, err := fh.bancoDados.BuscarFilmePorID(filmeID)
    if err != nil {
        // Fallback: retornar pelo menos o ID
        resposta := map[string]interface{}{
            "id":       filmeID,
            "mensagem": "Filme criado com sucesso",
        }
        enviarJSON(w, resposta, http.StatusCreated)
        return
    }
    
    fmt.Printf("✅ Filme criado: %s (ID: %d)\n", filmeCriado.Titulo, filmeCriado.ID)
    enviarJSON(w, filmeCriado, http.StatusCreated)
```

**Por que buscar novamente?**
- **Dados completos**: Cliente recebe filme com timestamps
- **Consistência**: Mostra exatamente como foi salvo
- **Status 201**: Indica criação bem-sucedida
- **Fallback**: Se busca falhar, pelo menos retorna ID

---

## 🚦 Sistema de Middleware

### Arquivo: `internal/handlers/middleware.go`

### O que é Middleware?

Middleware é um padrão que permite executar código **antes** e **depois** do handler principal:

```
Request → [Middleware 1] → [Middleware 2] → [Handler] → Response
            ↓                ↓              ↓
          Logs            CORS           Lógica
          Autenticação    Headers        Negócio
          Validação       Compressão     Database
```

### LogMiddleware

```go
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        inicio := time.Now()
        
        // Criar um ResponseWriter que captura o status code
        wrapperResposta := &responseWriter{
            ResponseWriter: w,
            statusCode:     http.StatusOK,
        }
        
        // Log da requisição
        fmt.Printf("🌐 %s %s - IP: %s - User-Agent: %s\n",
            r.Method,
            r.URL.Path,
            obterIPReal(r),
            r.UserAgent(),
        )
        
        // Executar próximo handler
        next.ServeHTTP(wrapperResposta, r)
        
        // Log da resposta
        duracao := time.Since(inicio)
        fmt.Printf("📊 %s %s - Status: %d - Duração: %v\n",
            r.Method,
            r.URL.Path,
            wrapperResposta.statusCode,
            duracao,
        )
    })
}
```

**Componentes do LogMiddleware:**

1. **Wrapper Function**: `func(next) func(w, r)`
2. **Timing**: `time.Now()` e `time.Since()`
3. **Response Wrapper**: Captura status code
4. **Informações úteis**: Método, URL, IP, User-Agent, duração

### ResponseWriter Wrapper

```go
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

**Por que wrapper?**
- `http.ResponseWriter` não expõe o status code escrito
- Precisamos "interceptar" a chamada `WriteHeader()`
- Salvar o código para usar no log

### CORSMiddleware

```go
func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Max-Age", "3600")
        
        // Se for preflight request, responder direto
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

**Headers CORS explicados:**
- **Allow-Origin**: Quais domínios podem fazer requests
- **Allow-Methods**: Que métodos HTTP são permitidos
- **Allow-Headers**: Que headers o cliente pode enviar
- **Max-Age**: Quanto tempo o browser pode cachear as regras CORS

### RecuperacaoMiddleware

```go
func RecuperacaoMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                fmt.Printf("🚨 PANIC recuperado: %v\n", err)
                
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusInternalServerError)
                
                resposta := `{"erro": "Erro interno do servidor", "codigo": 500}`
                w.Write([]byte(resposta))
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}
```

**Por que Recovery Middleware?**
- **Graceful Degradation**: Aplicação não quebra totalmente
- **Logs de Error**: Captura panics para debugging
- **Resposta Consistente**: Sempre retorna JSON estruturado
- **Disponibilidade**: Mantém servidor rodando

---

## 📊 Modelos de Dados Aprimorados

### Arquivo: `internal/models/filme.go` (Atualizado)

### Nova Struct: CriarFilme

```go
type CriarFilme struct {
    Titulo         string  `json:"titulo"`
    Descricao      string  `json:"descricao"`
    AnoLancamento  int     `json:"ano_lancamento"`
    DuracaoMinutos int     `json:"duracao_minutos"`
    Genero         string  `json:"genero"`
    Diretor        string  `json:"diretor"`
    Avaliacao      float64 `json:"avaliacao"`
}
```

**Por que struct separada para criação?**

**✅ Vantagens:**
- **Não inclui campos auto-gerados**: ID, DataCriacao, DataAtualizacao
- **Validação específica**: Regras diferentes para criação vs atualização
- **Flexibilidade**: Pode ter campos opcionais diferentes
- **Clareza**: Intent explícito (criar vs buscar)

**Comparação:**
```go
// ❌ Usando struct completa para criação
type Filme struct {
    ID              int       // Cliente não deve enviar
    DataCriacao     time.Time // Servidor controla
    DataAtualizacao time.Time // Servidor controla
    // ... outros campos
}

// ✅ Struct específica para criação
type CriarFilme struct {
    // Apenas campos que cliente deve fornecer
}
```

### RespostaErro Aprimorada

```go
type RespostaErro struct {
    Erro     string `json:"erro"`
    Codigo   int    `json:"codigo"`
    Detalhes string `json:"detalhes,omitempty"`
}
```

**Campo Detalhes:**
- **omitempty**: Só aparece se tiver conteúdo
- **Uso**: Erros de validação com múltiplos problemas
- **Exemplo**: "título é obrigatório; ano deve ser maior que 1888"

---

## 🔍 Sistema de Validação Robusto

### Arquivo: `internal/validators/filme_validator.go`

### Função ValidarCriarFilme

```go
func ValidarCriarFilme(filme *models.CriarFilme) []string {
    var erros []string
    
    // Validar título
    if strings.TrimSpace(filme.Titulo) == "" {
        erros = append(erros, "título é obrigatório")
    } else if len(filme.Titulo) > 255 {
        erros = append(erros, "título deve ter no máximo 255 caracteres")
    }
```

**Padrão de Validação:**
1. **Slice de erros**: Coleta todos os problemas
2. **Validação por campo**: Uma validação por vez
3. **Mensagens claras**: Usuário sabe exatamente o que corrigir
4. **strings.TrimSpace()**: Remove espaços antes/depois

### Validações Implementadas

#### 1. Título
```go
if strings.TrimSpace(filme.Titulo) == "" {
    erros = append(erros, "título é obrigatório")
} else if len(filme.Titulo) > 255 {
    erros = append(erros, "título deve ter no máximo 255 caracteres")
}
```

#### 2. Ano de Lançamento
```go
anoAtual := time.Now().Year()
if filme.AnoLancamento < 1888 { // Primeiro filme da história
    erros = append(erros, "ano de lançamento deve ser posterior a 1888")
} else if filme.AnoLancamento > anoAtual+5 { // Máximo 5 anos no futuro
    erros = append(erros, fmt.Sprintf("ano de lançamento não pode ser superior a %d", anoAtual+5))
}
```

**Por que 1888?**
- "Roundhay Garden Scene" (1888) é considerado o primeiro filme
- Validação baseada em fatos históricos

#### 3. Duração
```go
if filme.DuracaoMinutos < 1 {
    erros = append(erros, "duração deve ser maior que 0 minutos")
} else if filme.DuracaoMinutos > 600 { // 10 horas máximo
    erros = append(erros, "duração não pode exceder 600 minutos")
}
```

#### 4. Avaliação
```go
if filme.Avaliacao < 0 || filme.Avaliacao > 10 {
    erros = append(erros, "avaliação deve estar entre 0 e 10")
}
```

### Função LimparDados

```go
func LimparDados(filme *models.CriarFilme) {
    filme.Titulo = strings.TrimSpace(filme.Titulo)
    filme.Descricao = strings.TrimSpace(filme.Descricao)
    filme.Genero = strings.TrimSpace(filme.Genero)
    filme.Diretor = strings.TrimSpace(filme.Diretor)
}
```

**Sanitização de Dados:**
- Remove espaços extras
- Padroniza entrada do usuário
- Previne erros de validação por espaços

---

## 🗄️ Operações de Banco Atualizadas

### Arquivo: `internal/database/conexao.go` (Adição)

### Função CriarFilme

```go
func (bd *BancoDados) CriarFilme(filme *models.CriarFilme) (int, error) {
    query := `
        INSERT INTO filmes (titulo, descricao, ano_lancamento, duracao_minutos, genero, diretor, avaliacao)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `
    
    var novoID int
    
    err := bd.conexao.QueryRow(
        query,
        filme.Titulo,
        filme.Descricao,
        filme.AnoLancamento,
        filme.DuracaoMinutos,
        filme.Genero,
        filme.Diretor,
        filme.Avaliacao,
    ).Scan(&novoID)
    
    if err != nil {
        return 0, fmt.Errorf("erro ao inserir filme: %v", err)
    }
    
    fmt.Printf("💾 Filme inserido com ID: %d\n", novoID)
    return novoID, nil
}
```

**Componentes importantes:**

1. **INSERT com RETURNING**: PostgreSQL retorna o ID gerado
2. **Prepared Statement**: `$1, $2, ...` previne SQL injection
3. **QueryRow()**: Para comandos que retornam uma linha
4. **Scan()**: Captura o ID retornado
5. **Error Handling**: Contexto específico no erro

**RETURNING clause:**
```sql
-- ✅ PostgreSQL - retorna o ID criado
INSERT INTO filmes (...) VALUES (...) RETURNING id;

-- ❌ Alternativa menos elegante
INSERT INTO filmes (...) VALUES (...);
SELECT lastval(); -- Perigoso em ambiente concorrente
```

---

## 🚀 Main.go Simplificado e Organizado

### Arquivo: `cmd/server/main.go` (Novo)

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    
    "api-filmes/internal/database"
    "api-filmes/internal/handlers"
)

func main() {
    fmt.Println("🎬 Servidor da API de Filmes v2.0 iniciando...")
    
    // Conectar ao banco
    bancoDados, err := database.NovaConexao()
    if err != nil {
        log.Fatal("❌ Erro ao conectar com banco:", err)
    }
    
    defer func() {
        if err := bancoDados.Fechar(); err != nil {
            log.Printf("⚠️ Erro ao fechar conexão: %v", err)
        } else {
            fmt.Println("🔌 Conexão com banco fechada")
        }
    }()
    
    // Criar handler de filmes
    filmeHandler := handlers.NovoFilmeHandler(bancoDados)
    
    // Configurar rotas com middleware
    http.HandleFunc("/", aplicarMiddleware(filmeHandler.PaginaInicial))
    http.HandleFunc("/filmes", aplicarMiddleware(filmeHandler.ManipularFilmes))
    http.HandleFunc("/filmes/", aplicarMiddleware(filmeHandler.ManipularFilmeIndividual))
    
    // Iniciar servidor
    porta := ":8080"
    fmt.Printf("🚀 Servidor rodando em http://localhost%s\n", porta)
    
    if err := http.ListenAndServe(porta, nil); err != nil {
        log.Fatal("❌ Erro ao iniciar servidor:", err)
    }
}

// aplicarMiddleware combina todos os middlewares
func aplicarMiddleware(handler http.HandlerFunc) http.HandlerFunc {
    return handlers.RecuperacaoMiddleware(
        handlers.CORSMiddleware(
            handlers.LogMiddleware(handler),
        ),
    )
}
```

**Evolução do main.go:**

**✅ Agora:**
- **Apenas inicialização**: Foco na configuração
- **Dependency Injection**: Passa dependências para handlers
- **Middleware Chain**: Combina middlewares de forma elegante
- **Error Handling**: Tratamento adequado de erros de inicialização

**❌ Antes:**
- Misturava inicialização com lógica de negócio
- Variáveis globais
- Handlers soltos
- Sem middleware

### Cadeia de Middleware

```go
return handlers.RecuperacaoMiddleware(
    handlers.CORSMiddleware(
        handlers.LogMiddleware(handler),
    ),
)
```

**Ordem de execução:**
```
Request → Recuperacao → CORS → Log → Handler → Response
```

**Por que essa ordem?**
1. **Recuperacao**: Mais externo, captura qualquer panic
2. **CORS**: Headers precisam ser definidos cedo
3. **Log**: Registra após CORS, antes da lógica
4. **Handler**: Lógica de negócio por último

---

## 🧪 Testes Abrangentes

### 1. Preparação do Ambiente

```bash
# Verificar se banco está funcionando
go run cmd/server/main.go
```

**Saída esperada:**
```
🎬 Servidor da API de Filmes v2.0 iniciando...
🔌 Conectando ao banco de dados...
📍 Host: localhost:5432 | Banco: api_filmes
✅ Conexão com banco estabelecida com sucesso!
🚀 Servidor rodando em http://localhost:8080
```

### 2. Teste POST - Filme Válido

#### Configuração do Postman:
- **Método**: POST
- **URL**: `http://localhost:8080/filmes`
- **Headers**: 
  - `Content-Type: application/json`

#### Body (raw JSON):
```json
{
    "titulo": "Parasita",
    "descricao": "Uma família pobre se infiltra na casa de uma família rica",
    "ano_lancamento": 2019,
    "duracao_minutos": 132,
    "genero": "Thriller",
    "diretor": "Bong Joon-ho",
    "avaliacao": 8.6
}
```

#### Resposta Esperada (Status 201):
```json
{
    "id": 4,
    "titulo": "Parasita",
    "descricao": "Uma família pobre se infiltra na casa de uma família rica",
    "ano_lancamento": 2019,
    "duracao_minutos": 132,
    "genero": "Thriller",
    "diretor": "Bong Joon-ho",
    "avaliacao": 8.6,
    "data_criacao": "2024-01-20T15:30:45Z",
    "data_atualizacao": "2024-01-20T15:30:45Z"
}
```

#### Logs no Console:
```
🌐 POST /filmes - IP: 127.0.0.1:54321 - User-Agent: PostmanRuntime/7.32.2
➕ Criando novo filme...
💾 Filme inserido com ID: 4
✅ Filme criado: Parasita (ID: 4)
📊 POST /filmes - Status: 201 - Duração: 45ms
```

### 3. Teste POST - Dados Inválidos

#### Body com Múltiplos Erros:
```json
{
    "titulo": "",
    "descricao": "Descrição válida",
    "ano_lancamento": 1800,
    "duracao_minutos": -5,
    "genero": "",
    "diretor": "",
    "avaliacao": 15
}
```

#### Resposta Esperada (Status 400):
```json
{
    "erro": "Dados inválidos",
    "codigo": 400,
    "detalhes": "título é obrigatório; ano de lançamento deve ser posterior a 1888; duração deve ser maior que 0 minutos; gênero é obrigatório; diretor é obrigatório; avaliação deve estar entre 0 e 10"
}
```

### 4. Teste POST - Content-Type Incorreto

#### Headers:
- `Content-Type: text/plain`

#### Resposta Esperada (Status 400):
```json
{
    "erro": "Content-Type deve ser application/json",
    "codigo": 400
}
```

### 5. Teste POST - JSON Malformado

#### Body Inválido:
```json
{
    "titulo": "Filme"
    "ano": 2024    // Falta vírgula
}
```

#### Resposta Esperada (Status 400):
```json
{
    "erro": "JSON inválido",
    "codigo": 400
}
```

### 6. Teste POST - Campos Extras

#### Body com Campo Não Reconhecido:
```json
{
    "titulo": "Filme Teste",
    "ano_lancamento": 2024,
    "duracao_minutos": 120,
    "genero": "Ação",
    "diretor": "Diretor Teste",
    "avaliacao": 8.0,
    "campo_extra": "valor não permitido"
}
```

#### Resposta Esperada (Status 400):
```json
{
    "erro": "JSON inválido",
    "codigo": 400
}
```

### 7. Teste GET - Verificar Filme Criado

#### Requisição:
- **Método**: GET
- **URL**: `http://localhost:8080/filmes/4`

#### Resposta:
```json
{
    "id": 4,
    "titulo": "Parasita",
    // ... dados completos
}
```

### 8. Teste GET - Lista Atualizada

#### Requisição:
- **Método**: GET
- **URL**: `http://localhost:8080/filmes`

#### Verificar:
- Total aumentou
- Novo filme aparece na lista
- Ordenação mantida

---

## 🎓 Conceitos Aprendidos

### 1. Organização de Código
- **Handlers em Structs**: Organização por funcionalidade
- **Dependency Injection**: Passagem de dependências
- **Separation of Concerns**: Cada arquivo com responsabilidade específica

### 2. Middleware Pattern
- **Chain of Responsibility**: Cada middleware tem uma responsabilidade
- **Cross-cutting Concerns**: Logs, CORS, Recovery aplicados a todas as rotas
- **Composição**: Combinação de múltiplos middlewares

### 3. Request Processing
- **Content-Type Validation**: Garantir formato correto
- **JSON Decoding**: Conversão de JSON para struct
- **Data Validation**: Verificação de regras de negócio
- **Error Handling**: Diferentes tipos de erro com status codes apropriados

### 4. HTTP Status Codes
- **200 OK**: Operação bem-sucedida (GET)
- **201 Created**: Recurso criado com sucesso (POST)
- **400 Bad Request**: Dados inválidos do cliente
- **409 Conflict**: Violação de restrição (título duplicado)
- **500 Internal Server Error**: Erro no servidor

### 5. Validation Strategies
- **Multiple Validation**: Coletar todos os erros de uma vez
- **Business Rules**: Validações baseadas em regras de negócio
- **Data Sanitization**: Limpeza de dados antes da validação
- **Contextual Messages**: Mensagens específicas para cada erro

---

## 🏗️ Padrões de Arquitetura Implementados

### 1. Handler Pattern
```go
// ✅ Organizado em struct
type FilmeHandler struct {
    bancoDados *database.BancoDados
}

func (fh *FilmeHandler) criarFilme(w http.ResponseWriter, r *http.Request) {
    // Lógica específica
}

// ❌ Funções globais desordenadas
func criarFilme(w http.ResponseWriter, r *http.Request) {
    // usa variável global
}
```

### 2. Middleware Chain Pattern
```go
// Middleware como Higher-Order Function
func LogMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Lógica antes
        next.ServeHTTP(w, r)
        // Lógica depois
    })
}

// Composição de middlewares
aplicarMiddleware(handler) // Recovery -> CORS -> Log -> Handler
```

### 3. Data Transfer Object (DTO) Pattern
```go
// Entrada específica para criação
type CriarFilme struct {
    Titulo string `json:"titulo"`
    // Apenas campos necessários para criação
}

// Entidade completa para resposta
type Filme struct {
    ID              int       `json:"id"`
    DataCriacao     time.Time `json:"data_criacao"`
    // Todos os campos incluindo auto-gerados
}
```

### 4. Repository Pattern (Implícito)
```go
// Interface implícita
type FilmeRepository interface {
    CriarFilme(*models.CriarFilme) (int, error)
    BuscarFilmePorID(int) (*models.Filme, error)
    // Operações de dados abstraídas
}

// Implementação concreta
type BancoDados struct {
    conexao *sql.DB
}
```

### 5. Error Wrapping Pattern
```go
// Adicionar contexto aos erros
if err != nil {
    return 0, fmt.Errorf("erro ao inserir filme: %v", err)
}

// Tratamento específico por tipo
if strings.Contains(err.Error(), "duplicate") {
    enviarErro(w, "Filme já existe", http.StatusConflict)
}
```

---

## 🔄 Fluxo de Dados Detalhado

### Request → Response Flow (POST /filmes)

```
1. Cliente (Postman/Frontend)
   ↓ POST /filmes + JSON body
   
2. Go HTTP Server
   ↓ http.ListenAndServe
   
3. Middleware Chain
   ↓ Recovery → CORS → Log
   
4. FilmeHandler.ManipularFilmes()
   ↓ switch r.Method = "POST"
   
5. FilmeHandler.criarFilme()
   ↓ Content-Type validation
   
6. JSON Decoding
   ↓ json.NewDecoder(r.Body).Decode()
   
7. Data Validation
   ↓ validators.ValidarCriarFilme()
   
8. Database Layer
   ↓ bancoDados.CriarFilme()
   
9. PostgreSQL
   ↓ INSERT ... RETURNING id
   
10. Response Construction
    ↓ Buscar filme completo + JSON encoding
    
11. HTTP Response
    ↓ Status 201 + headers + JSON body
    
12. Cliente recebe filme criado
```

### Error Flow Examples

```
Erro de Validação:
Client → JSON inválido → Validation → 400 Bad Request

Erro de Duplicação:
Client → Título existe → PostgreSQL constraint → 409 Conflict

Erro de Servidor:
Client → Dados válidos → Database down → 500 Internal Server Error
```

---

## 🛡️ Aspectos de Segurança Avançados

### 1. Input Validation
```go
// ✅ Validação robusta
if strings.TrimSpace(filme.Titulo) == "" {
    erros = append(erros, "título é obrigatório")
}

// ✅ Sanitização
filme.Titulo = strings.TrimSpace(filme.Titulo)

// ✅ Limite de tamanho
if len(filme.Titulo) > 255 {
    erros = append(erros, "título muito longo")
}
```

### 2. JSON Security
```go
// ✅ Rejeitar campos desconhecidos
decoder.DisallowUnknownFields()

// ✅ Verificar Content-Type
if r.Header.Get("Content-Type") != "application/json" {
    return BadRequest
}

// ✅ Limitar tamanho do body (implementar quando necessário)
r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB
```

### 3. SQL Injection Prevention
```go
// ✅ SEGURO - Prepared statements
query := "INSERT INTO filmes (titulo) VALUES ($1)"
db.QueryRow(query, filme.Titulo)

// ❌ VULNERÁVEL - String concatenation
// query := "INSERT INTO filmes (titulo) VALUES ('" + filme.Titulo + "')"
```

### 4. Error Information Disclosure
```go
// ✅ SEGURO - Log detalhado interno, mensagem genérica externa
fmt.Printf("❌ Erro específico: %v\n", err)           // Log
enviarErro(w, "Erro interno", 500)                   // Cliente

// ❌ PERIGOSO - Vazar detalhes internos
// enviarErro(w, err.Error(), 500)  // Pode expor paths, senhas, etc.
```

### 5. CORS Configuration
```go
// 🚨 Desenvolvimento - permite tudo
w.Header().Set("Access-Control-Allow-Origin", "*")

// ✅ Produção - específico
// w.Header().Set("Access-Control-Allow-Origin", "https://meusite.com")
```

---

## 📊 Monitoramento e Observabilidade

### Logs Estruturados
```go
// Request logs
fmt.Printf("🌐 %s %s - IP: %s - User-Agent: %s\n",
    r.Method, r.URL.Path, obterIPReal(r), r.UserAgent())

// Operation logs  
fmt.Printf("➕ Criando novo filme...\n")
fmt.Printf("💾 Filme inserido com ID: %d\n", novoID)

// Performance logs
fmt.Printf("📊 %s %s - Status: %d - Duração: %v\n",
    r.Method, r.URL.Path, status, duracao)

// Error logs
fmt.Printf("❌ Erro ao salvar filme: %v\n", err)
```

### Métricas Básicas
```go
// Timing
inicio := time.Now()
// ... operação ...
duracao := time.Since(inicio)

// Status tracking
wrapperResposta.statusCode // Capturado pelo middleware

// Request info
r.Method, r.URL.Path, r.UserAgent(), obterIPReal(r)
```

### Health Indicators
```go
// Database connectivity
if err := conexao.Ping(); err != nil {
    return nil, fmt.Errorf("erro ao conectar com banco: %v", err)
}

// Successful operations
fmt.Printf("✅ Filme criado: %s (ID: %d)\n", titulo, id)
```

---

## 🔧 Troubleshooting Avançado

### Problema: "JSON inválido" mas JSON parece correto
```bash
# Verificar encoding
file -I arquivo.json

# Verificar caracteres especiais
hexdump -C arquivo.json | head

# Usar validator online
# https://jsonlint.com
```

### Problema: "Content-Type deve ser application/json"
```bash
# No Postman, verificar:
# 1. Headers tab → Content-Type: application/json
# 2. Body tab → raw → JSON (não Text)
# 3. Não ter trailing spaces

# Via curl:
curl -X POST http://localhost:8080/filmes \
  -H "Content-Type: application/json" \
  -d '{"titulo": "Teste"}'
```

### Problema: Middleware não executa
```go
// ✅ Verificar se aplicarMiddleware está sendo usado
http.HandleFunc("/filmes", aplicarMiddleware(handler.ManipularFilmes))

// ❌ Sem middleware
// http.HandleFunc("/filmes", handler.ManipularFilmes)

// ✅ Verificar ordem dos middlewares
aplicarMiddleware(handler) // Recovery → CORS → Log → Handler
```

### Problema: "Filme já existe" mas não deveria
```sql
-- Verificar constraint de unicidade
\d filmes

-- Ver se existe index único no título
\di

-- Se não houver, criar:
-- CREATE UNIQUE INDEX idx_filmes_titulo ON filmes(titulo);
```

### Problema: Status code sempre 200
```go
// ❌ WriteHeader depois de Write
w.Write([]byte("dados"))
w.WriteHeader(400)  // Muito tarde!

// ✅ WriteHeader antes de Write
w.WriteHeader(400)
w.Write([]byte("dados"))

// ✅ Ou usar helper
enviarJSON(w, dados, 400)  // Faz na ordem correta
```

---

## 📈 Performance e Otimização

### Database Operations
```go
// ✅ Usar prepared statements (já implementado)
db.QueryRow("SELECT * FROM filmes WHERE id = $1", id)

// ✅ Buscar apenas campos necessários
SELECT id, titulo, diretor FROM filmes  // Para listagem
SELECT * FROM filmes WHERE id = $1      // Para detalhes

// 🔮 Para otimizações futuras:
// - Connection pooling configuration
// - Query caching
// - Database indexes
```

### JSON Processing
```go
// ✅ Stream processing (já implementado)
json.NewDecoder(r.Body).Decode(&struct)  // Não carrega tudo na memória

// ✅ Struct tags otimizadas
type Filme struct {
    ID int `json:"id"`  // Simples e direto
}

// 🔮 Para otimizações futuras:
// - JSON streaming para listas grandes
// - Compression middleware
```

### Memory Management
```go
// ✅ Resource cleanup (já implementado)
defer linhas.Close()
defer bancoDados.Fechar()

// ✅ Response writer wrapper eficiente
type responseWriter struct {
    http.ResponseWriter
    statusCode int  // Apenas o necessário
}
```

---

## 🎯 Preparação para Módulo 3-B

### O que já funciona perfeitamente:
- ✅ Criação de filmes via POST
- ✅ Validação robusta de dados
- ✅ Middleware funcionando
- ✅ Logs detalhados
- ✅ Tratamento de erros contextual
- ✅ Organização clara do código

### O que vamos adicionar no Módulo 3-B:
- 🔜 **PUT /filmes/{id}** - Atualizar filme completo
- 🔜 **PATCH /filmes/{id}** - Atualizar campos específicos
- 🔜 **DELETE /filmes/{id}** - Remover filme
- 🔜 **Validações específicas** para atualização
- 🔜 **Concorrência básica** com timestamps
- 🔜 **Soft delete** vs hard delete

### Conceitos que aprenderemos:
- **HTTP PUT vs PATCH**: Diferenças conceituais e práticas
- **Partial Updates**: Como atualizar apenas alguns campos
- **Optimistic Locking**: Prevenção de atualizações concorrentes
- **Audit Trail**: Rastreamento de mudanças
- **Cascade Operations**: Operações relacionadas

---

## 🚀 Comparação: Módulo 2 vs Módulo 3-A

### Módulo 2 (Apenas Leitura):
```go
// Estrutura simples
func listarFilmes(w http.ResponseWriter, r *http.Request) {
    filmes, _ := banco.BuscarTodos()
    json.NewEncoder(w).Encode(filmes)
}

// Sem validação
// Sem middleware
// Código no main.go
```

### Módulo 3-A (Leitura + Escrita):
```go
// Estrutura organizada
type FilmeHandler struct {
    bancoDados *database.BancoDados
}

func (fh *FilmeHandler) criarFilme(w http.ResponseWriter, r *http.Request) {
    // 1. Validar Content-Type
    // 2. Decodificar JSON
    // 3. Validar dados
    // 4. Salvar no banco
    // 5. Retornar resultado
}

// Com middleware completo
// Com validação robusta
// Código organizado em pacotes
```

**Evolução alcançada:**
- 🎯 **Arquitetura profissional** vs código iniciante
- 🛡️ **Segurança robusta** vs básica
- 📊 **Observabilidade completa** vs logs mínimos
- 🔧 **Manutenibilidade alta** vs código monolítico
- 🚀 **Escalabilidade preparada** vs estrutura simples

---

## 📚 Referências e Estudos Adicionais

### Documentação Oficial
- [HTTP Package](https://pkg.go.dev/net/http) - Servidor HTTP
- [JSON Package](https://pkg.go.dev/encoding/json) - Manipulação JSON
- [Strings Package](https://pkg.go.dev/strings) - Manipulação de strings
- [Time Package](https://pkg.go.dev/time) - Operações com tempo

### Padrões e Práticas
- [Effective Go](https://golang.org/doc/effective_go.html) - Boas práticas
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) - Convenções
- [REST API Guidelines](https://github.com/microsoft/api-guidelines) - Padrões REST

### Segurança
- [OWASP Go Secure Coding](https://github.com/OWASP/Go-SCP) - Segurança em Go
- [Input Validation](https://cheatsheetseries.owasp.org/cheatsheets/Input_Validation_Cheat_Sheet.html) - Validação

### Ferramentas de Desenvolvimento
- **Postman**: Teste de APIs
- **JSON Validator**: Validação de JSON online
- **DB Browser**: Visualização de dados
- **Go Playground**: Teste de código Go

---

## ✅ Checklist Final do Módulo 3-A

### Configuração e Estrutura:
- [ ] Nova estrutura de pastas criada
- [ ] Handlers organizados em struct
- [ ] Middleware funcionando
- [ ] Validadores implementados

### Funcionalidade POST:
- [ ] POST /filmes cria filme com dados válidos
- [ ] Retorna status 201 Created
- [ ] Valida Content-Type application/json
- [ ] Rejeita campos desconhecidos
- [ ] Valida todos os campos obrigatórios
- [ ] Retorna erros de validação detalhados

### Logs e Monitoramento:
- [ ] Logs de requisição aparecem
- [ ] Logs de operação aparecem
- [ ] Logs de performance aparecem
- [ ] Logs de erro aparecem

### Testes Realizados:
- [ ] Filme válido criado com sucesso
- [ ] Dados inválidos rejeitados com 400
- [ ] Content-Type incorreto rejeitado
- [ ] JSON malformado rejeitado
- [ ] Campos extras rejeitados
- [ ] Filme aparece na listagem GET

### Compreensão:
- [ ] Entendo diferença entre GET e POST
- [ ] Sei como funciona middleware chain
- [ ] Compreendo validação de dados
- [ ] Reconheço padrões de arquitetura aplicados
- [ ] Entendo fluxo request → response completo

---

**🎉 Parabéns! Você completou o Módulo 3-A e agora tem uma API que não apenas lê dados, mas também cria novos registros de forma segura e organizada!**

**No Módulo 3-B, vamos completar o CRUD implementando as operações de atualização e exclusão, tornando nossa API totalmente funcional!**