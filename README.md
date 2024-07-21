# Desafio Go: Client-Server para Cotação de Dólar

Este projeto consiste em dois sistemas desenvolvidos em Go: `client.go` e `server.go`. O objetivo é criar uma comunicação entre cliente e servidor para obter a cotação do dólar e registrar as cotações em um banco de dados SQLite.

## Requisitos

1. `client.go` deverá fazer uma requisição HTTP para `server.go` solicitando a cotação do dólar.
2. `server.go` deverá consumir a API de câmbio de Dólar e Real no endereço: `https://economia.awesomeapi.com.br/json/last/USD-BRL`.
3. `server.go` deverá retornar o resultado no formato JSON para o cliente.
4. `server.go` deverá registrar no banco de dados SQLite cada cotação recebida.
5. O tempo máximo para chamar a API de cotação do dólar deverá ser de 200ms.
6. O tempo máximo para persistir os dados no banco deverá ser de 10ms.
7. `client.go` deverá receber do `server.go` apenas o valor atual do câmbio (campo "bid" do JSON).
8. O tempo máximo para `client.go` receber o resultado do `server.go` deverá ser de 300ms.
9. Os três contextos (requisição da API, persistência no banco e comunicação entre client-server) deverão retornar erro nos logs caso o tempo de execução seja insuficiente.
10. `client.go` deverá salvar a cotação atual em um arquivo `cotacao.txt` no formato: `Dólar: {valor}`.
11. O endpoint necessário gerado pelo `server.go` será `/cotacao` e a porta utilizada pelo servidor HTTP será a 8080.

## Instruções

### 1. Clonar o Repositório

```bash
git clone https://github.com/pimentafm/client-server-go.git
cd client-server-go
```
### 2. Executar o server

```bash
cd server
go run server.go
```
### 3. Executar o client

```bash
cd client
go run client.go
```

