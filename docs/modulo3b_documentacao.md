# Módulo 3-B: Operações de Atualização e Exclusão (PUT/DELETE)
## 🔄 Completando o CRUD da API de Filmes

### 📖 Objetivos do Módulo
- Implementar endpoint PUT para atualizar filmes completos
- Implementar endpoint DELETE para remover filmes
- Adicionar validações específicas para operações de atualização
- Implementar verificação de existência de recursos
- Criar endpoint de estatísticas para insights dos dados
- Completar o CRUD (Create, Read, Update, Delete) da API

---

## 🧠 Conceitos Fundamentais

### O que é CRUD?
CRUD é um acrônimo que representa as quatro operações básicas de persistência de dados:

- **C**reate (Criar) → POST
- **R**ead (Ler) → GET
- **U**pdate (Atualizar) → PUT/PATCH
- **D**elete (Deletar) → DELETE

### HTTP Methods para CRUD

| Operação | Método HTTP | Endpoint | Propósito |
|----------|-------------|----------|-----------|
| Create | POST | /filmes | Criar novo filme |
| Read (All) | GET | /filmes | Listar todos os filmes |
| Read (One) | GET | /filmes/{id} | Buscar filme específico |
| Update | PUT | /filmes/{id} | Atualizar filme completo |
| Delete | DELETE | /filmes/{id} | Remover filme |

### PUT vs PATCH
```
PUT - Substituição Completa:
- Substitui o recurso inteiro
- Todos os campos devem ser fornecidos
- Idempotente (mesmo resultado se chamado múltiplas vezes)

PATCH - Atualização Parcial:
- Atualiza apenas campos específicos
- Permite atualizações incrementais
- Mais complexo de implementar
```

### Resource Existence Checking
Antes de atualizar ou deletar, sempre verificar se o recurso existe:

```
1. Cliente solicita PUT /filmes/123
2. Servidor verifica se filme 123 existe
3. Se não existir → 404 Not Found
4. Se existir → procede com atualização
```

---

## 📊 Evolução dos Modelos de Dados

### Arquivo: `internal/models/filme.go` (Atualizado)

### Nova Struct: AtualizarFilme

```go
type AtualizarFilme struct {
    Titulo         string  `json:"titulo"`
    Descricao      string  `json:"descricao"`
    AnoLancamento  int     `json:"ano_lancamento"`
    DuracaoMinutos int     `json:"duracao_minutos"`
    Genero         string  `json:"genero"`
    Diretor        string  `json:"diretor"`
    Avaliacao      float64 `json:"avaliacao"`
}
```

**Por que struct separada para atualização?**

**Comparação com CriarFilme:**
```go
// CriarFilme e AtualizarFilme são idênticas atualmente
// Mas podem divergir no futuro:

type CriarFilme struct {
    // Campos obrigatórios na criação
    Titulo string `json:"titulo" validate:"required"`
}

type AtualizarFilme struct {
    // Talvez permita título opcional no futuro
    Titulo *string `json:"titulo,omitempty"`
}
```

**Benefícios:**
- **Flexibilidade futura**: Regras diferentes para criar vs atualizar
- **Validações específicas**: Campos obrigatórios podem diferir
- **Evolução independente**: Cada operação pode evoluir separadamente
- **Clareza de intent**: Código mais legível e autodocumentado

### Nova Struct: EstatisticasFilmes

```go
type EstatisticasFilmes struct {
    TotalFilmes        int     `json:"total_filmes"`
    AvaliacaoMedia     float64 `json:"avaliacao_media"`
    DuracaoMediaMinutos float64 `json:"duracao_media_minutos"`
    GeneroMaisComum    string  `json:"genero_mais_comum"`
}
```

**Propósito:**
- **Business Intelligence**: Insights sobre os dados
- **Dashboard**: Informações para interfaces administrativas
- **Monitoramento**: Acompanhar crescimento da base de dados
- **Relatórios**: Dados agregados para relatórios

---

## 🔍 Sistema de Validação Expandido

### Arquivo: `internal/validators/filme_validator.go` (Atualizado)

### Função ValidarAtualizarFilme

```go
func ValidarAtualizarFilme(filme *models.AtualizarFilme) []string {
    var erros []string
    
    // Validações similares ao criar, mas com algumas diferenças
    if strings.TrimSpace(filme.Titulo) == "" {
        erros = append(erros, "título é obrigatório")
    } else if len(filme.Titulo) > 255 {
        erros = append(erros, "título deve ter no máximo 255 caracteres")
    }
    
    // ... outras validações
    
    return erros
}
```

**Por que função separada?**

1. **Regras específicas**: Atualização pode ter regras diferentes
2. **Campos opcionais**: No futuro, alguns campos podem ser opcionais na atualização
3. **Validações contextuais**: Pode precisar do estado atual do recurso
4. **Flexibilidade**: Permite evolução independente das validações

### Função LimparDadosAtualizar

```go
func LimparDadosAtualizar(filme *models.AtualizarFilme) {
    filme.Titulo = strings.TrimSpace(filme.Titulo)
    filme.Descricao = strings.TrimSpace(filme.Descricao)
    filme.Genero = strings.TrimSpace(filme.Genero)
    filme.Diretor = strings.TrimSpace(filme.Diretor)
}
```

**Sanitização consistente:**
- Remove espaços em branco desnecessários
- Padroniza dados antes da validação
- Previne erros por diferenças de formatação

### Função ValidarExistenciaFilme

```go
func ValidarExistenciaFilme(filmeID int, bancoDados interface{ BuscarFilmePorID(int) (*models.Filme, error) }) error {
    _, err := bancoDados.BuscarFilmePorID(filmeID)
    if err != nil {
        if strings.Contains(err.Error(), "não encontrado") {
            return fmt.Errorf("filme com ID %d não encontrado", filmeID)
        }
        return fmt.Errorf("erro ao verificar existência do filme: %v", err)
    }
    return nil
}
```

**Interface Duck Typing:**
- Usa interface implícita para flexibilidade
- Permite mocking em testes
- Não depende de tipos concretos

---

## 🗄️ Operações de Banco de Dados Avançadas

### Arquivo: `internal/database/conexao.go` (Adições)

### Função AtualizarFilme

```go
func (bd *BancoDados) AtualizarFilme(id int, filme *models.AtualizarFilme) error {
    query := `
        UPDATE filmes 
        SET titulo = $1, 
            descricao = $2, 
            ano_lancamento = $3, 
            duracao_minutos = $4, 
            genero = $5, 
            diretor = $6, 
            avaliacao = $7,
            data_atualizacao = CURRENT_TIMESTAMP
        WHERE id = $8
    `
    
    resultado, err := bd.conexao.Exec(
        query,
        filme.Titulo,
        filme.Descricao,
        filme.AnoLancamento,
        filme.DuracaoMinutos,
        filme.Genero,
        filme.Diretor,
        filme.Avaliacao,
        id,
    )
    
    if err != nil {
        return fmt.Errorf("erro ao atualizar filme: %v", err)
    }
    
    // Verificar se alguma linha foi afetada
    linhasAfetadas, err := resultado.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    if linhasAfetadas == 0 {
        return fmt.Errorf("filme com ID %d não encontrado", id)
    }
    
    fmt.Printf("🔄 Filme ID %d atualizado com sucesso\n", id)
    return nil
}
```

**Componentes importantes:**

1. **UPDATE com WHERE**: Atualiza apenas o registro específico
2. **CURRENT_TIMESTAMP**: PostgreSQL atualiza automaticamente o timestamp
3. **Exec() vs QueryRow()**: Exec para comandos que não retornam dados
4. **RowsAffected()**: Verifica se a operação teve impacto
5. **Error Context**: Adiciona contexto específico aos erros

**Por que verificar RowsAffected?**
```sql
-- Se não existir filme com ID 999:
UPDATE filmes SET titulo = 'Novo' WHERE id = 999;
-- Query executa sem erro, mas não afeta nenhuma linha
-- RowsAffected() = 0 → sabemos que não encontrou o registro
```

### Função DeletarFilme

```go
func (bd *BancoDados) DeletarFilme(id int) error {
    query := "DELETE FROM filmes WHERE id = $1"
    
    resultado, err := bd.conexao.Exec(query, id)
    if err != nil {
        return fmt.Errorf("erro ao deletar filme: %v", err)
    }
    
    // Verificar se alguma linha foi afetada
    linhasAfetadas, err := resultado.RowsAffected()
    if err != nil {
        return fmt.Errorf("erro ao verificar linhas afetadas: %v", err)
    }
    
    if linhasAfetadas == 0 {
        return fmt.Errorf("filme com ID %d não encontrado", id)
    }
    
    fmt.Printf("🗑️ Filme ID %d deletado com sucesso\n", id)
    return nil
}
```

**Hard Delete vs Soft Delete:**

```sql
-- Hard Delete (implementado)
DELETE FROM filmes WHERE id = $1;

-- Soft Delete (alternativa)
UPDATE filmes SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1;
```

**Quando usar cada um:**
- **Hard Delete**: Dados realmente devem sumir (LGPD, cleanup)
- **Soft Delete**: Auditoria, recovery, histórico

### Função VerificarExistenciaFilme

```go
func (bd *BancoDados) VerificarExistenciaFilme(id int) (bool, error) {
    query := "SELECT EXISTS(SELECT 1 FROM filmes WHERE id = $1)"
    
    var existe bool
    err := bd.conexao.QueryRow(query, id).Scan(&existe)
    if err != nil {
        return false, fmt.Errorf("erro ao verificar existência do filme: %v", err)
    }
    
    return existe, nil
}
```

**Por que usar EXISTS?**
```sql
-- ✅ Eficiente - para apenas na primeira correspondência
SELECT EXISTS(SELECT 1 FROM filmes WHERE id = $1);

-- ❌ Menos eficiente - conta todos os registros
SELECT COUNT(*) FROM filmes WHERE id = $1;

-- ❌ Traz dados desnecessários
SELECT * FROM filmes WHERE id = $1;
```

### Função ObterEstatisticasFilmes

```go
func (bd *BancoDados) ObterEstatisticasFilmes() (*models.EstatisticasFilmes, error) {
    stats := &models.EstatisticasFilmes{}
    
    // Query para obter estatísticas básicas
    query := `
        SELECT 
            COUNT(*) as total,
            ROUND(AVG(avaliacao), 2) as avaliacao_media,
            ROUND(AVG(duracao_minutos), 2) as duracao_media
        FROM filmes
    `
    
    err := bd.conexao.QueryRow(query).Scan(
        &stats.TotalFilmes,
        &stats.AvaliacaoMedia,
        &stats.DuracaoMediaMinutos,
    )
    
    if err != nil {
        return nil, fmt.Errorf("erro ao obter estatísticas básicas: %v", err)
    }
    
    // Query para obter gênero mais comum
    queryGenero := `
        SELECT genero 
        FROM filmes 
        GROUP BY genero 
        ORDER BY COUNT(*) DESC 
        LIMIT 1
    `
    
    err = bd.conexao.QueryRow(queryGenero).Scan(&stats.GeneroMaisComum)
    if err != nil {
        // Se não houver filmes, definir valores padrão
        stats.GeneroMaisComum = "N/A"
    }
    
    return stats, nil
}
```

**Funções SQL utilizadas:**
- **COUNT(*)**: Conta total de registros
- **AVG()**: Calcula média aritmética
- **ROUND(valor, casas)**: Arredonda para número específico de casas decimais
- **GROUP BY**: Agrupa registros por critério
- **ORDER BY COUNT(*) DESC**: Ordena por frequência (maior primeiro)
- **LIMIT 1**: Retorna apenas o primeiro resultado

**Tratamento de casos edge:**
```go
// Se tabela vazia, gênero fica como "N/A"
if err != nil {
    stats.GeneroMaisComum = "N/A"
}
```

---

## 🎮 Handlers Avançados

### Arquivo: `internal/handlers/filme_handlers.go` (Atualizações)

### ManipularFilmeIndividual Atualizado

```go
func (fh *FilmeHandler) ManipularFilmeIndividual(w http.ResponseWriter, r *http.Request) {
    configurarCabecalhos(w)
    
    // Extrair ID da URL
    caminho := strings.TrimPrefix(r.URL.Path, "/filmes/")
    if caminho == "" {
        enviarErro(w, "ID do filme é obrigatório", http.StatusBadRequest)
        return
    }
    
    // Verificar se é rota de estatísticas
    if caminho == "estatisticas" {
        if r.Method == "GET" {
            fh.obterEstatisticas(w, r)
        } else {
            enviarErro(w, "Método não permitido", http.StatusMethodNotAllowed)
        }
        return
    }
    
    id, err := strconv.Atoi(caminho)
    if err != nil {
        enviarErro(w, "ID inválido", http.StatusBadRequest)
        return
    }
    
    switch r.Method {
    case "GET":
        fh.buscarFilmePorID(w, r, id)
    case "PUT":
        fh.atualizarFilme(w, r, id)  // ✨ NOVO
    case "DELETE":
        fh.deletarFilme(w, r, id)   // ✨ NOVO
    case "OPTIONS":
        w.WriteHeader(http.StatusOK)
    default:
        enviarErro(w, "Método não permitido", http.StatusMethodNotAllowed)
    }
}
```

**Roteamento inteligente:**
1. **Extração de path**: Remove prefixo `/filmes/`
2. **Rota especial**: `/filmes/estatisticas` tem tratamento específico
3. **Validação de ID**: Converte para inteiro e valida
4. **Method routing**: Distribui para função específica baseado no método HTTP

### Método atualizarFilme

```go
func (fh *FilmeHandler) atualizarFilme(w http.ResponseWriter, r *http.Request, id int) {
    fmt.Printf("🔄 Atualizando filme ID: %d\n", id)
    
    // Verificar Content-Type
    if r.Header.Get("Content-Type") != "application/json" {
        enviarErro(w, "Content-Type deve ser application/json", http.StatusBadRequest)
        return
    }
    
    // Verificar se filme existe antes de tentar atualizar
    existe, err := fh.bancoDados.VerificarExistenciaFilme(id)
    if err != nil {
        fmt.Printf("❌ Erro ao verificar existência do filme: %v\n", err)
        enviarErro(w, "Erro interno do servidor", http.StatusInternalServerError)
        return
    }
    
    if !existe {
        enviarErro(w, fmt.Sprintf("Filme com ID %d não encontrado", id), http.StatusNotFound)
        return
    }
    
    // ... resto da implementação
}
```

**Fluxo de validação:**
1. **Content-Type**: Garante que é JSON
2. **Existence check**: Verifica se recurso existe
3. **Early return**: Falha rápida se não passar nas validações
4. **Specific errors**: Mensagens específicas para cada tipo de erro

**Por que verificar existência primeiro?**
- **User Experience**: Erro 404 imediato vs processamento desnecessário
- **Performance**: Evita parsing JSON se recurso não existir
- **Security**: Não vaza informações sobre IDs válidos vs inválidos

### Método deletarFilme

```go
func (fh *FilmeHandler) deletarFilme(w http.ResponseWriter, r *http.Request, id int) {
    fmt.Printf("🗑️ Deletando filme ID: %d\n", id)
    
    // Verificar se filme existe antes de tentar deletar
    filme, err := fh.bancoDados.BuscarFilmePorID(id)
    if err != nil {
        if strings.Contains(err.Error(), "não encontrado") {
            enviarErro(w, fmt.Sprintf("Filme com ID %d não encontrado", id), http.StatusNotFound)
        } else {
            fmt.Printf("❌ Erro ao verificar filme: %v\n", err)
            enviarErro(w, "Erro interno do servidor", http.StatusInternalServerError)
        }
        return
    }
    
    // Salvar informações para o log antes de deletar
    tituloFilme := filme.Titulo
    
    // Deletar do banco
    err = fh.bancoDados.DeletarFilme(id)
    if err != nil {
        fmt.Printf("❌ Erro ao deletar filme: %v\n", err)
        
        if strings.Contains(err.Error(), "não encontrado") {
            enviarErro(w, fmt.Sprintf("Filme com ID %d não encontrado", id), http.StatusNotFound)
        } else {
            enviarErro(w, "Erro interno do servidor", http.StatusInternalServerError)
        }
        return
    }
    
    fmt.Printf("✅ Filme deletado: %s (ID: %d)\n", tituloFilme, id)
    
    // Retornar confirmação
    resposta := models.RespostaSucesso{
        Mensagem: fmt.Sprintf("Filme '%s' deletado com sucesso", tituloFilme),
        Dados: map[string]interface{}{
            "id":     id,
            "titulo": tituloFilme,
        },
    }
    
    enviarJSON(w, resposta, http.StatusOK)
}
```

**Estratégia de delete com informações:**
1. **Buscar primeiro**: Obtém dados antes de deletar
2. **Salvar informações**: Título para logs e resposta
3. **Delete operation**: Remove do banco
4. **Informative response**: Retorna dados sobre o que foi deletado

**Por que não usar apenas DeletarFilme()?**
```go
// ❌ Menos informativo
err := bancoDados.DeletarFilme(id)
resposta := "Filme deletado"

// ✅ Mais informativo
filme := bancoDados.BuscarFilmePorID(id)  // Salva info
err := bancoDados.DeletarFilme(id)
resposta := fmt.Sprintf("Filme '%s' deletado", filme.Titulo)
```

### Método obterEstatisticas

```go
func (fh *FilmeHandler) obterEstatisticas(w http.ResponseWriter, r *http.Request) {
    fmt.Println("📊 Obtendo estatísticas dos filmes...")
    
    stats, err := fh.bancoDados.ObterEstatisticasFilmes()
    if err != nil {
        fmt.Printf("❌ Erro ao obter estatísticas: %v\n", err)
        enviarErro(w, "Erro interno do servidor", http.StatusInternalServerError)
        return
    }
    
    fmt.Printf("✅ Estatísticas obtidas: %d filmes total\n", stats.TotalFilmes)
    enviarJSON(w, stats, http.StatusOK)
}
```

**Endpoint de agregação:**
- **Business Intelligence**: Dados para dashboards
- **Simple implementation**: Operação direta de leitura
- **Cacheable**: Pode ser otimizado com cache no futuro

---

## 🧪 Testes Abrangentes

### 1. Configuração do Ambiente

```bash
# Verificar estrutura atualizada
ls -la internal/
# Deve mostrar: handlers/ models/ database/ config/ validators/

# Executar servidor
go run cmd/server/main.go
```

**Saída esperada:**
```
🎬 Servidor da API de Filmes v3.0 iniciando...
🔌 Conectando ao banco de dados...
📍 Host: localhost:5432 | Banco: api_filmes
✅ Conexão com banco estabelecida com sucesso!
🚀 Servidor rodando em http://localhost:8080
```

### 2. Teste PUT - Atualização de Filme

#### Configuração Postman:
- **Método**: PUT
- **URL**: `http://localhost:8080/filmes/1`
- **Headers**: `Content-Type: application/json`

#### Body (Atualização Completa):
```json
{
    "titulo": "O Poderoso Chefão - Edição Especial",
    "descricao": "A saga completa de uma família mafiosa italiana nos Estados Unidos - versão remasterizada com cenas inéditas",
    "ano_lancamento": 1972,
    "duracao_minutos": 180,
    "genero": "Drama/Crime",
    "diretor": "Francis Ford Coppola",
    "avaliacao": 9.5
}
```

#### Resposta Esperada (Status 200):
```json
{
    "id": 1,
    "titulo": "O Poderoso Chefão - Edição Especial",
    "descricao": "A saga completa de uma família mafiosa italiana nos Estados Unidos - versão remasterizada com cenas inéditas",
    "ano_lancamento": 1972,
    "duracao_minutos": 180,
    "genero": "Drama/Crime",
    "diretor": "Francis Ford Coppola",
    "avaliacao": 9.5,
    "data_criacao": "2024-01-15T10:30:00Z",
    "data_atualizacao": "2024-01-20T16:45:30Z"
}
```

#### Logs no Console:
```
🌐 PUT /filmes/1 - IP: 127.0.0.1:54321 - User-Agent: PostmanRuntime/7.32.2
🔄 Atualizando filme ID: 1
🔄 Filme ID 1 atualizado com sucesso
✅ Filme atualizado: O Poderoso Chefão - Edição Especial (ID: 1)
📊 PUT /filmes/1 - Status: 200 - Duração: 52ms
```

### 3. Teste PUT - Filme Inexistente

#### Requisição:
- **URL**: `http://localhost:8080/filmes/999`
- **Body**: Qualquer JSON válido

#### Resposta Esperada (Status 404):
```json
{
    "erro": "Filme com ID 999 não encontrado",
    "codigo": 404
}
```

### 4. Teste PUT - Dados Inválidos

#### Body com Erros:
```json
{
    "titulo": "",
    "descricao": "Descrição válida",
    "ano_lancamento": 1800,
    "duracao_minutos": -10,
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

### 5. Teste DELETE - Remoção de Filme

#### Requisição:
- **Método**: DELETE
- **URL**: `http://localhost:8080/filmes/3`

#### Resposta Esperada (Status 200):
```json
{
    "mensagem": "Filme 'Vingadores: Ultimato' deletado com sucesso",
    "dados": {
        "id": 3,
        "titulo": "Vingadores: Ultimato"
    }
}
```

#### Logs no Console:
```
🌐 DELETE /filmes/3 - IP: 127.0.0.1:54321 - User-Agent: PostmanRuntime/7.32.2
🗑️ Deletando filme ID: 3
🗑️ Filme ID 3 deletado com sucesso
✅ Filme deletado: Vingadores: Ultimato (ID: 3)
📊 DELETE /filmes/3 - Status: 200 - Duração: 38ms
```

### 6. Teste DELETE - Filme Inexistente

#### Requisição:
- **URL**: `http://localhost:8080/filmes/999`

#### Resposta Esperada (Status 404):
```json
{
    "erro": "Filme com ID 999 não encontrado",
    "codigo": 404
}
```

### 7. Teste GET - Verificar Lista Atualizada

#### Após Updates e Deletes:
- **URL**: `http://localhost:8080/filmes`

#### Verificar:
- Total de filmes diminuiu (após deletes)
- Filmes atualizados mostram novos dados
- `data_atualizacao` diferente de `data_criacao` nos atualizados

### 8. Teste GET - Estatísticas

#### Requisição:
- **Método**: GET
- **URL**: `http://localhost:8080/filmes/estatisticas`

#### Resposta Esperada:
```json
{
    "total_filmes": 3,
    "avaliacao_media": 8.73,
    "duracao_media_minutos": 152.33,
    "genero_mais_comum": "Drama"
}
```

#### Logs no Console:
```
🌐 GET /filmes/estatisticas - IP: 127.0.0.1:54321 - User-Agent: PostmanRuntime/7.32.2
📊 Obtendo estatísticas dos filmes...
✅ Estatísticas obtidas: 3 filmes total
📊 GET /filmes/estatisticas - Status: 200 - Duração: 15ms
```

### 9. Teste de Sequência Completa

#### Workflow CRUD Completo:
1. **POST** `/filmes` - Criar filme
2. **GET** `/filmes/{id}` - Verificar criação
3. **PUT** `/filmes/{id}` - Atualizar filme
4. **GET** `/filmes/{id}` - Verificar atualização
5. **DELETE** `/filmes/{id}` - Deletar filme
6. **GET** `/filmes/{id}` - Verificar deleção (404)
7. **GET** `/filmes/estatisticas` - Ver estatísticas atualizadas

---

## 🎓 Conceitos Aprendidos

### 1. CRUD Operations
- **Complete lifecycle**: Criação → Leitura → Atualização → Deleção
- **HTTP method mapping**: POST/GET/PUT/DELETE
- **Resource identification**: URLs com IDs
- **Status code usage**: 200, 201, 404, 400, 500

### 2. Resource Management
- **Existence checking**: Verificar antes de operar
- **Atomic operations**: Uma operação por request
- **Idempotency**: PUT pode ser chamado múltiplas vezes
- **Resource state**: Timestamps para auditoria

### 3. Database Operations
- **UPDATE with WHERE**: Modificação específica
- **DELETE with WHERE**: Remoção específica
- **RowsAffected()**: Verificação de impacto
- **Aggregate queries**: COUNT, AVG, GROUP BY

### 4. Error Handling Strategies
- **Preemptive validation**: Verificar antes de processar
- **Specific error codes**: 404 vs 400 vs 500
- **Contextual messages**: Erros específicos para cada situação
- **Graceful degradation**: Fallbacks quando possível

### 5. API Design Patterns
- **RESTful URLs**: `/recursos/{id}` para operações específicas
- **HTTP semantics**: Métodos apropriados para cada operação
- **Consistent responses**: Formato padrão para sucessos e erros
- **Resource representations**: DTOs específicos para cada operação

---

## 🏗️ Padrões de Arquitetura Consolidados

### 1. Repository Pattern (Completo)
```go
// Interface implícita completa
type FilmeRepository interface {
    // Create
    CriarFilme(*models.CriarFilme) (int, error)
    
    // Read
    BuscarTodosFilmes() ([]models.FilmeResumo, error)
    BuscarFilmePorID(int) (*models.Filme, error)
    ContarFilmes() (int, error)
    
    // Update
    AtualizarFilme(int, *models.AtualizarFilme) error
    
    // Delete
    DeletarFilme(int) error
    
    // Utility
    VerificarExistenciaFilme(int) (bool, error)
    ObterEstatisticasFilmes() (*models.EstatisticasFilmes, error)
}

// BancoDados implementa implicitamente FilmeRepository
type BancoDados struct {
    conexao *sql.DB
}
```

### 2. Handler Pattern (Organizado)
```go
// Handlers agrupados por recurso
type FilmeHandler struct {
    bancoDados *database.BancoDados
}

// Método para cada operação HTTP
func (fh *FilmeHandler) ManipularFilmes(w, r)      // POST, GET /filmes
func (fh *FilmeHandler) ManipularFilmeIndividual(w, r) // GET, PUT, DELETE /filmes/{id}
func (fh *FilmeHandler) obterEstatisticas(w, r)    // GET /filmes/estatisticas
```

### 3. DTO Pattern (Data Transfer Objects)
```go
// Diferentes DTOs para diferentes operações
type CriarFilme struct { /* campos para criação */ }
type AtualizarFilme struct { /* campos para atualização */ }
type FilmeResumo struct { /* campos para listagem */ }
type Filme struct { /* entidade completa */ }
type EstatisticasFilmes struct { /* dados agregados */ }
```

### 4. Validation Strategy Pattern
```go
// Validações específicas por operação
func ValidarCriarFilme(*models.CriarFilme) []string
func ValidarAtualizarFilme(*models.AtualizarFilme) []string
func ValidarID(int) error
func ValidarExistenciaFilme(int, repository) error
```

### 5. Error Handling Strategy
```go
// Tratamento em camadas
Database Layer  → Erros técnicos com contexto
Handler Layer   → Tradução para HTTP status codes
Response Layer  → Mensagens amigáveis para cliente
```

---

## 🔄 Fluxo de Dados Completo (PUT/DELETE)

### PUT /filmes/{id} Flow

```
1. Cliente (Postman/Frontend)
   ↓ PUT /filmes/1 + JSON body
   
2. Go HTTP Server
   ↓ http.ListenAndServe + middleware chain
   
3. FilmeHandler.ManipularFilmeIndividual()
   ↓ extract ID + route to PUT handler
   
4. FilmeHandler.atualizarFilme()
   ↓ validate Content-Type
   
5. Existence Check
   ↓ bancoDados.VerificarExistenciaFilme(1)
   
6. JSON Decoding & Validation
   ↓ json.Decode + validators.ValidarAtualizarFilme()
   
7. Database Update
   ↓ bancoDados.AtualizarFilme(1, dados)
   
8. PostgreSQL
   ↓ UPDATE filmes SET ... WHERE id = $1
   
9. Fetch Updated Resource
   ↓ bancoDados.BuscarFilmePorID(1)
   
10. HTTP Response
    ↓ Status 200 + JSON do filme atualizado
    
11. Cliente recebe filme atualizado
```

### DELETE /filmes/{id} Flow

```
1. Cliente (Postman/Frontend)
   ↓ DELETE /filmes/3
   
2. Go HTTP Server
   ↓ http.ListenAndServe + middleware chain
   
3. FilmeHandler.ManipularFilmeIndividual()
   ↓ extract ID + route to DELETE handler
   
4. FilmeHandler.deletarFilme()
   ↓ fetch film info for logging
   
5. Existence Check & Info Retrieval
   ↓ bancoDados.BuscarFilmePorID(3)
   
6. Save Film Info
   ↓ titulo := filme.Titulo (for response)
   
7. Database Delete
   ↓ bancoDados.DeletarFilme(3)
   
8. PostgreSQL
   ↓ DELETE FROM filmes WHERE id = $1
   
9. HTTP Response
   ↓ Status 200 + confirmation message
   
10. Cliente recebe confirmação de deleção
```

### Error Flows

```
PUT Error Flow:
Client → Dados inválidos → Validation → 400 Bad Request
Client → ID inexistente → Existence check → 404 Not Found
Client → Título duplicado → Database constraint → 409 Conflict

DELETE Error Flow:
Client → ID inexistente → Existence check → 404 Not Found
Client → Constraint violation → Database error → 500 Internal Error
```

---

## 🛡️ Aspectos de Segurança Avançados

### 1. Input Validation Robusta
```go
// ✅ Validação em múltiplas camadas
1. Content-Type validation
2. JSON schema validation (DisallowUnknownFields)
3. Business rules validation
4. Database constraint validation

// ✅ Sanitização consistente
func LimparDadosAtualizar(filme *models.AtualizarFilme) {
    filme.Titulo = strings.TrimSpace(filme.Titulo)
    // Remove caracteres especiais se necessário
}
```

### 2. Resource Authorization (Para futuro)
```go
// 🔮 Preparado para autorização
func (fh *FilmeHandler) atualizarFilme(w, r, id) {
    // Futuro: verificar se usuário pode editar este filme
    // if !userCanEdit(userID, filmID) { return 403 }
    
    // Atual: qualquer um pode editar qualquer filme
}
```

### 3. SQL Injection Prevention
```go
// ✅ SEGURO - Prepared statements em todas as operações
UPDATE filmes SET titulo = $1 WHERE id = $2
DELETE FROM filmes WHERE id = $1
SELECT EXISTS(SELECT 1 FROM filmes WHERE id = $1)

// ❌ VULNERÁVEL (não usado)
// "UPDATE filmes SET titulo = '" + titulo + "' WHERE id = " + id
```

### 4. Information Disclosure Prevention
```go
// ✅ SEGURO - Não vazar detalhes internos
if err != nil {
    log.Printf("Internal error: %v", err)           // Para developer
    enviarErro(w, "Erro interno", 500)              // Para cliente
}

// ❌ PERIGOSO
// enviarErro(w, err.Error(), 500)  // Pode vazar paths, senhas, etc.
```

### 5. Resource Enumeration
```go
// 🚨 ATENÇÃO - Permite enumerar IDs válidos
// GET /filmes/1 → 200 (existe)
// GET /filmes/999 → 404 (não existe)
// Atacante pode descobrir quais IDs existem

// 🔮 Para mitigar no futuro:
// - Rate limiting por IP
// - Autenticação obrigatória
// - UUIDs em vez de IDs sequenciais
```

---

## 📊 Performance e Otimização

### Database Operations
```go
// ✅ Otimizações já implementadas
1. Prepared statements (precompiladas)
2. Specific field selection (SELECT id, titulo vs SELECT *)
3. EXISTS() vs COUNT(*) para existence checks
4. Connection pooling automático (Go database/sql)

// 🔮 Otimizações futuras
1. Database indexes nos campos mais consultados
2. Query caching para estatísticas
3. Batch operations para múltiplas atualizações
4. Read replicas para operações de leitura
```

### Memory Management
```go
// ✅ Gestão eficiente de recursos
defer linhas.Close()           // Libera cursors
defer bancoDados.Fechar()      // Libera conexões

// ✅ Streaming JSON
json.NewDecoder(r.Body)        // Não carrega JSON inteiro na memória
json.NewEncoder(w)             // Stream direto para response
```

### HTTP Performance
```go
// ✅ Headers otimizados
Content-Type: application/json; charset=utf-8
Access-Control-Max-Age: 3600   // Cache preflight por 1 hora

// 🔮 Melhorias futuras
// - Compression middleware (gzip)
// - ETag headers para caching
// - HTTP/2 support
```

---

## 🔧 Troubleshooting Avançado

### Problema: PUT retorna 404 mas recurso existe
```bash
# Verificar se ID está sendo extraído corretamente
curl -v PUT http://localhost:8080/filmes/1
# Logs devem mostrar: "🔄 Atualizando filme ID: 1"

# Se aparecer ID diferente, problema na extração da URL
# Verificar: strings.TrimPrefix(r.URL.Path, "/filmes/")
```

### Problema: DELETE aparenta funcionar mas não remove do banco
```sql
-- Verificar se UPDATE está sendo usado em vez de DELETE
SELECT * FROM filmes WHERE id = 123;

-- Se filme ainda existe, verificar:
-- 1. Query SQL está correta?
-- 2. Transação sendo commitada?
-- 3. Conexão está estável?
```

### Problema: RowsAffected sempre retorna 0
```go
// ✅ Verificar se está usando Exec() corretamente
resultado, err := bd.conexao.Exec(query, params...)
if err != nil {
    return err  // Erro na query
}

linhas, err := resultado.RowsAffected()
if err != nil {
    return err  // Erro ao obter affected rows
}

// ❌ Não usar Query() para UPDATE/DELETE
// rows, err := bd.conexao.Query(query, params...)  // ERRADO
```

### Problema: Estatísticas retornam valores incorretos
```sql
-- Verificar dados na base
SELECT COUNT(*), AVG(avaliacao), AVG(duracao_minutos) FROM filmes;

-- Verificar se ROUND está funcionando
SELECT ROUND(8.333333, 2);  -- Deve retornar 8.33

-- Verificar GROUP BY para gênero mais comum
SELECT genero, COUNT(*) FROM filmes GROUP BY genero ORDER BY COUNT(*) DESC;
```

### Problema: Middleware não executa em algumas rotas
```go
// ✅ Verificar se todas as rotas usam aplicarMiddleware
http.HandleFunc("/filmes", aplicarMiddleware(handler.ManipularFilmes))
http.HandleFunc("/filmes/", aplicarMiddleware(handler.ManipularFilmeIndividual))

// ❌ Rota sem middleware
// http.HandleFunc("/filmes", handler.ManipularFilmes)  // SEM LOGS
```

### Problema: CORS ainda bloqueia requests
```bash
# Verificar headers de resposta
curl -I http://localhost:8080/filmes
# Deve incluir:
# Access-Control-Allow-Origin: *
# Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS

# Se não aparecer, verificar se CORSMiddleware está ativo
```

---

## 📈 Monitoramento e Observabilidade

### Logs Estruturados por Operação

#### PUT Logs:
```bash
🌐 PUT /filmes/1 - IP: 127.0.0.1 - User-Agent: PostmanRuntime/7.32.2
🔄 Atualizando filme ID: 1
🔄 Filme ID 1 atualizado com sucesso
✅ Filme atualizado: Novo Título (ID: 1)
📊 PUT /filmes/1 - Status: 200 - Duração: 52ms
```

#### DELETE Logs:
```bash
🌐 DELETE /filmes/3 - IP: 127.0.0.1 - User-Agent: PostmanRuntime/7.32.2
🗑️ Deletando filme ID: 3
🗑️ Filme ID 3 deletado com sucesso
✅ Filme deletado: Vingadores: Ultimato (ID: 3)
📊 DELETE /filmes/3 - Status: 200 - Duração: 38ms
```

#### Error Logs:
```bash
🌐 PUT /filmes/999 - IP: 127.0.0.1 - User-Agent: PostmanRuntime/7.32.2
🔄 Atualizando filme ID: 999
❌ Erro ao verificar existência do filme: filme com ID 999 não encontrado
📊 PUT /filmes/999 - Status: 404 - Duração: 15ms
```

### Métricas Importantes
```go
// Request metrics
Method, URL, StatusCode, Duration, IP, UserAgent

// Business metrics  
FilmesCreated, FilmesUpdated, FilmesDeleted

// Performance metrics
DatabaseQueryTime, JSONProcessingTime, ValidationTime

// Error metrics
404Count, 500Count, ValidationErrors
```

### Health Checks
```go
// Database connectivity
func (bd *BancoDados) HealthCheck() error {
    return bd.conexao.Ping()
}

// Statistics as health indicator
stats, err := bancoDados.ObterEstatisticasFilmes()
// Se TotalFilmes == 0, pode indicar problema
```

---

## 🚀 Comparação: Evolução Completa dos Módulos

### Módulo 1 (Básico):
```go
// Apenas GET com dados estáticos
func listarFilmes(w, r) {
    json := `{"filmes": [...]}`
    w.Write([]byte(json))
}

// Estrutura: Apenas main.go
// Funcionalidade: Leitura de dados fixos
// Tecnologia: HTTP básico + strings
```

### Módulo 2 (Com Banco):
```go
// GET com dados dinâmicos
func listarFilmes(w, r) {
    filmes, _ := banco.BuscarTodos()
    json.NewEncoder(w).Encode(filmes)
}

// Estrutura: main.go + database + models + config
// Funcionalidade: Leitura de banco real
// Tecnologia: PostgreSQL + structs + JSON
```

### Módulo 3-A (Com POST):
```go
// GET + POST com validação
type FilmeHandler struct { banco *database.BancoDados }
func (fh *FilmeHandler) criarFilme(w, r) {
    // Validação + criação + middleware
}

// Estrutura: handlers + validators + middleware
// Funcionalidade: Leitura + criação com validação
// Tecnologia: Arquitetura organizada + middleware
```

### Módulo 3-B (CRUD Completo):
```go
// GET + POST + PUT + DELETE + estatísticas
func (fh *FilmeHandler) ManipularFilmeIndividual(w, r) {
    switch r.Method {
    case "GET": fh.buscarFilmePorID(w, r, id)
    case "PUT": fh.atualizarFilme(w, r, id)
    case "DELETE": fh.deletarFilme(w, r, id)
    }
}

// Estrutura: Arquitetura completa e profissional
// Funcionalidade: CRUD completo + estatísticas
// Tecnologia: API REST completa + padrões de arquitetura
```

### Evolução Quantitativa:

| Aspecto | Módulo 1 | Módulo 2 | Módulo 3-A | Módulo 3-B |
|---------|----------|----------|------------|------------|
| **Arquivos** | 1 | 4 | 7 | 8 |
| **Linhas de código** | ~50 | ~200 | ~400 | ~600 |
| **Endpoints** | 2 | 3 | 4 | 7 |
| **HTTP Methods** | GET | GET | GET, POST | GET, POST, PUT, DELETE |
| **Validações** | 0 | Básica | Robusta | Completa |
| **Middleware** | 0 | 0 | 3 tipos | 3 tipos |
| **Tratamento erro** | Básico | Médio | Avançado | Profissional |

---

## 🎯 Preparação para Próximos Passos

### O que temos agora (API Completa):
- ✅ **CRUD completo**: Create, Read, Update, Delete
- ✅ **Validação robusta**: Dados sempre consistentes
- ✅ **Middleware profissional**: Logs, CORS, Recovery
- ✅ **Arquitetura escalável**: Código organizado e manutenível
- ✅ **Tratamento de erros**: Status codes apropriados
- ✅ **Estatísticas**: Insights sobre os dados
- ✅ **Documentação**: Guias completos de cada módulo

### Próximas evoluções possíveis:

#### 🔄 **Funcionalidades Avançadas:**
- **Paginação**: `GET /filmes?page=1&limit=10`
- **Filtros**: `GET /filmes?genero=Drama&ano=2019`
- **Ordenação**: `GET /filmes?sort=avaliacao&order=desc`
- **Busca**: `GET /filmes/search?q=padrinho`

#### 🛡️ **Segurança:**
- **Autenticação**: JWT tokens
- **Autorização**: Roles e permissões
- **Rate Limiting**: Limite de requests por IP
- **Input Sanitization**: XSS prevention

#### 📊 **Performance:**
- **Cache**: Redis para dados frequentes
- **Database Indexes**: Otimização de queries
- **Connection Pooling**: Configuração avançada
- **Compression**: Gzip middleware

#### 🧪 **Qualidade:**
- **Testes unitários**: Cobertura de código
- **Testes integração**: Testes end-to-end
- **Linting**: Code quality tools
- **CI/CD**: Pipelines automatizados

#### 🚀 **Deployment:**
- **Dockerização**: Containers
- **Kubernetes**: Orquestração
- **Monitoring**: Prometheus + Grafana
- **Logging**: Structured logs + ELK Stack

#### 📈 **Escalabilidade:**
- **Microservices**: Separação por domínio
- **Message Queues**: Processamento assíncrono
- **Load Balancing**: Múltiplas instâncias
- **Database Sharding**: Distribuição de dados

---

## 📚 Referências e Estudos Avançados

### Documentação Oficial
- [HTTP Package](https://pkg.go.dev/net/http) - Servidor HTTP completo
- [SQL Package](https://pkg.go.dev/database/sql) - Interface de banco
- [JSON Package](https://pkg.go.dev/encoding/json) - Serialização
- [PostgreSQL Docs](https://www.postgresql.org/docs/) - Database específico

### Padrões e Arquitetura
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - Robert Martin
- [REST API Guidelines](https://github.com/microsoft/api-guidelines) - Microsoft
- [Go Project Layout](https://github.com/golang-standards/project-layout) - Padrões Go

### Segurança
- [OWASP API Security](https://owasp.org/www-project-api-security/) - Top 10 vulnerabilidades
- [Go Security](https://github.com/OWASP/Go-SCP) - Secure coding practices

### Performance
- [High Performance Go](https://dave.cheney.net/high-performance-go-workshop/dotgo-paris.html) - Dave Cheney
- [Go Database Patterns](https://www.alexedwards.net/blog/organising-database-access) - Alex Edwards

### Ferramentas
- **Testing**: [Testify](https://github.com/stretchr/testify) - Framework de testes
- **Mocking**: [GoMock](https://github.com/golang/mock) - Mock generation
- **Documentation**: [Swagger](https://swagger.io/) - API documentation
- **Monitoring**: [Prometheus](https://prometheus.io/) - Metrics collection

---

## ✅ Checklist Final do Módulo 3-B

### Implementação CRUD:
- [ ] PUT /filmes/{id} atualiza filme existente
- [ ] PUT retorna 404 para filme inexistente
- [ ] PUT valida dados antes de atualizar
- [ ] PUT retorna filme atualizado completo
- [ ] DELETE /filmes/{id} remove filme existente
- [ ] DELETE retorna 404 para filme inexistente
- [ ] DELETE retorna confirmação com dados do filme removido

### Funcionalidades Adicionais:
- [ ] GET /filmes/estatisticas retorna dados agregados
- [ ] Verificação de existência antes de operações
- [ ] Logs detalhados para todas as operações
- [ ] Tratamento específico de erros por operação

### Validações e Segurança:
- [ ] Validação de dados em atualizações
- [ ] Prepared statements em todas as queries
- [ ] Headers CORS configurados
- [ ] Recovery middleware captura panics

### Testes Realizados:
- [ ] PUT com dados válidos (sucesso)
- [ ] PUT com dados inválidos (400)
- [ ] PUT com ID inexistente (404)
- [ ] DELETE com ID válido (sucesso)
- [ ] DELETE com ID inexistente (404)
- [ ] GET estatísticas (dados corretos)
- [ ] Workflow CRUD completo testado

### Compreensão:
- [ ] Entendo diferenças entre GET/POST/PUT/DELETE
- [ ] Sei implementar CRUD completo
- [ ] Compreendo padrões de arquitetura aplicados
- [ ] Reconheço importância de validação de existência
- [ ] Entendo fluxo completo de request → response
- [ ] Sei como fazer troubleshooting de problemas

---

**🎉 PARABÉNS! Você completou todo o sistema CRUD e agora tem uma API REST profissional e completa!**

**Sua jornada de aprendizado em Go e APIs REST alcançou um marco importante. Você construiu uma aplicação real, funcional e seguindo as melhores práticas da indústria!**

**O próximo passo é escolher quais funcionalidades avançadas implementar para levar sua API ao próximo nível! 🚀**