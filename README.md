# Projeto de API com Go e PostgreSQL

## Introdução

Este projeto é uma API simples escrita em Go que usa PostgreSQL como banco de dados. O objetivo é fornecer endpoints para criar, ler, atualizar e excluir tweets.

## Tecnologias Utilizadas

- Go
- PostgreSQL
- Podman

## Instalação

### Pré-requisitos

- Go 1.20+
- PostgreSQL
- Podman

### Passos para Instalação

1. Clone o Repositório:

    ```bash
    git clone https://github.com/seuusuario/seurepositorio.git
    cd seurepositorio
    ```

2. Configurar o Banco de Dados:

    ```sql
    CREATE DATABASE seu_banco;
    \c seu_banco;
    CREATE TABLE tweets (
        id SERIAL PRIMARY KEY,
        description TEXT NOT NULL
    );
    ```

3. Configurar o Podman:

    ```bash
    podman run -d --name seu_postgres -e POSTGRES_USER=user -e POSTGRES_PASSWORD=toor -e POSTGRES_DB=db -p 5432:5432 postgres
    ```

4. Compilar e Rodar o Projeto:

    ```bash
    go build -o seuapp
    ./seuapp
    ```

## Uso

### Endpoints da API

- **GET /tweets**: Obtém todos os tweets.
- **POST /tweets**: Cria um novo tweet.
- **PUT /tweets/:id**: Atualiza um tweet existente.
- **DELETE /tweets/:id**: Exclui um tweet existente.

### Exemplos de Requisições

- **GET Todos os Tweets**:

    ```bash
    curl -X GET http://localhost:8080/tweets
    ```

- **POST Criar Tweet**:

    ```bash
    curl -X POST http://localhost:8080/tweets \
      -H "Content-Type: application/json" \
      -d '{"description":"Novo tweet"}'
    ```

- **PUT Atualizar Tweet**:

    ```bash
    curl -X PUT http://localhost:8080/tweets/1 \
      -H "Content-Type: application/json" \
      -d '{"description":"Tweet atualizado"}'
    ```

- **DELETE Excluir Tweet**:

    ```bash
    curl -X DELETE http://localhost:8080/tweets/1
    ```

## Estrutura do Projeto

- `main.go`: Arquivo principal que inicializa o servidor.
- `controllers/`: Contém os handlers da API.
- `db/`: Contém a configuração do banco de dados.
- `entities/`: Contém as definições de modelos.
- `routes/`: Contém as rotas da api.
