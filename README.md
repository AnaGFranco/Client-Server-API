### README.md

# Pós Go Expert - Desafio: Client-Server-API

## Como Executar

### Pré-requisitos

- [Go](https://golang.org/doc/install) (versão 1.16 ou superior)
- [SQLite](https://www.sqlite.org/download.html)

### Passo a Passo

1. Clone o repositório:
   ```sh
   git clone git@github.com:AnaGFranco/Client-Server-API.git
   cd Client-Server-API
   ```

2. Inicialize o banco de dados:
    - O banco de dados será criado automaticamente ao rodar o `server/main.go`.


3. Execute o servidor:
   ```sh
   cd server
   go run main.go
   ```
    - O servidor estará ouvindo na porta `8080`.
   

4. Em outro terminal, execute o cliente:
   ```sh
   cd client
   go run main.go
   ```
    - O cliente fará uma requisição ao servidor para obter a cotação do dólar e salvará o valor no arquivo `cotacao.txt`.
