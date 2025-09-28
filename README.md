# 🎬 API de Filmes - Portfolio

> API REST completa desenvolvida em Go com PostgreSQL, containerizada com Docker

[![Go](https://img.shields.io/badge/Go-1.21-blue.svg)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-blue.svg)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-blue.svg)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

🌐 API: http://localhost:8081  # Em vez de 8080

## 📋 Sobre o Projeto

Uma API REST completa para gerenciamento de filmes, desenvolvida seguindo as melhores práticas de desenvolvimento em Go. O projeto demonstra conhecimentos em:

- **Backend Development** com Go
- **Database Design** com PostgreSQL
- **RESTful API** design
- **Docker** containerization
- **Clean Architecture** principles

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
git clone https://github.com/seu-usuario/api-filmes.git
cd api-filmes

# Inicie a aplicação
make run

# Ou usando docker-compose diretamente
docker-compose up -d