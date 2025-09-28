# 🎬 API de Filmes - Portfolio

> API REST completa desenvolvida em Go com PostgreSQL, containerizada com Docker

[![Go](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-blue.svg)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

🌐 API: http://localhost:8081  # Em vez de 8080

## 📋 Sobre o Projeto

Uma API REST completa para gerenciamento de filmes, desenvolvida seguindo as melhores práticas de desenvolvimento em Go. O objetivo principal deste repositório é compartilhar conhecimento e apoiar quem deseja compreender como construir APIs REST do zero, com uma abordagem prática, simples e didática. Além de entregar uma API funcional e containerizada, o projeto serve como material de estudo com uma trilha de módulos que explica decisões, padrões e boas práticas passo a passo.

- **Backend Development** com Go
- **Database Design** com PostgreSQL
- **RESTful API** design
- **Docker** containerization
- **Clean Architecture** principles

## 📘 Aprendizado por Módulos (docs/)
Este repositório inclui uma trilha de aprendizado dentro da pasta `docs/`. Você pode seguir os módulos na ordem para aprender desde os conceitos iniciais até funcionalidades mais avançadas da API. Essa trilha é ideal para quem está começando ou migrando para Go, e quer ver exemplos práticos.

- Módulo 1: [docs/modulo1_documentacao.md](docs/modulo1_documentacao.md)
- Módulo 2: [docs/modulo2_documentacao.md](docs/modulo2_documentacao.md)
- Módulo 3A: [docs/modulo3a_documentacao.md](docs/modulo3a_documentacao.md)
- Módulo 3B: [docs/modulo3b_documentacao.md](docs/modulo3b_documentacao.md)
- Módulo 4A (Final): [docs/modulo4a_final_documentacao.md](docs/modulo4a_final_documentacao.md)
- Módulo 4B (Final): [docs/modulo4b_final_documentacao.md](docs/modulo4b_final_documentacao.md)

Sugestão de leitura:
1. Comece pelos módulos 1 e 2 para entender a base do projeto.
2. Siga pelos módulos 3A e 3B conforme evolui a API.
3. Finalize com os módulos 4A e 4B, consolidando o conhecimento e revisando boas práticas.

Cada módulo aprofunda o entendimento da arquitetura, dos handlers HTTP, do acesso a dados (PostgreSQL), da validação e do empacotamento com Docker.

## 📘 Aprendizado por Módulos (docs/)
Este repositório inclui uma trilha de aprendizado passo a passo dentro da pasta `docs/`. Você pode seguir os módulos na ordem para compreender desde os conceitos iniciais até funcionalidades mais avançadas da API.

- Módulo 1: [docs/modulo1_documentacao.md](docs/modulo1_documentacao.md)
- Módulo 2: [docs/modulo2_documentacao.md](docs/modulo2_documentacao.md)
- Módulo 3A: [docs/modulo3a_documentacao.md](docs/modulo3a_documentacao.md)
- Módulo 3B: [docs/modulo3b_documentacao.md](docs/modulo3b_documentacao.md)
- Módulo 4A (Final): [docs/modulo4a_final_documentacao.md](docs/modulo4a_final_documentacao.md)
- Módulo 4B (Final): [docs/modulo4b_final_documentacao.md](docs/modulo4b_final_documentacao.md)

Sugestão de leitura:
1. Comece pelos módulos 1 e 2 para entender a base do projeto.
2. Siga pelos módulos 3A e 3B conforme evolui a API.
3. Finalize com os módulos 4A e 4B, consolidando o conhecimento e revisando boas práticas.

Cada módulo aprofunda o entendimento da arquitetura, dos handlers HTTP, do acesso a dados (PostgreSQL), da validação e do empacotamento com Docker.

### 🛠️ Tecnologias Utilizadas

- **[Go 1.21](https://golang.org/)** - Linguagem de programação
- **[PostgreSQL 15](https://www.postgresql.org/)** - Banco de dados
- **[Docker](https://www.docker.com/)** - Containerização
- **[Docker Compose](https://docs.docker.com/compose/)** - Orquestração

## 🚀 Quick Start

### Pré-requisitos

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Executando o Projeto
```bash
# Clone o repositório
git clone https://github.com/nabazio/api-filmes.git
cd api-filmes

# Inicie a aplicação
make run

# Ou usando docker-compose diretamente
docker-compose up -d
```
## 🤝 Contribuição e Feedback
Este projeto existe para ensinar e compartilhar. Sugestões, dúvidas e melhorias são muito bem-vindas! 
- Abra uma issue com seu questionamento ou ideia
- Envie um PR com ajustes ou novas seções de documentação
- Compartilhe o repositório com quem está aprendendo APIs em Go

Se este conteúdo te ajudar, deixe uma ⭐ no repositório para apoiar a iniciativa.