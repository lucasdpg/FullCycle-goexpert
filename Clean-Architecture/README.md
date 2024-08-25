# Clean Architecture

## Desafio proposto.

Olá devs!
Agora é a hora de botar a mão na massa. Para este desafio, você precisará criar o usecase de listagem das orders.
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.


### Acessando o projeto.
`git clone git@github.com:lucasdpg/FullCycle-goexpert.git && cd FullCycle-goexpert/Clean-Architecture`

### Iniciando Mysql e RabbitMQ.
`docker-compose up -d`
Obs: Pode ser necessário esperar alguns minutos ate que o banco fique pronto para rodar o próximo passo.

### Preparar o Banco de Dados
`make migrate`
Obs: Em alguns momentos este comando pode falahar a conexnão com o banco, nestes casos reiniciei o banco de dados com `docker-compose down && docker-compose up -d`

### Start do projeto.
`cd cmd/ordersystem && go run main.go wire_gen.go`

### Informações uteis para validar o desafio:

Se tudo ocorreu como deveria os serviços devem subir em WEB, gRPC e GraphQL segue as portas:

WEB: 8000, gRPC: 50051 e GraphQL: 8080

Os serviços devem ser acessados em `localhost`

### WEBSERVER
Acessando os arquivos api/create_order.http e api/list_order.http encontrara as chamadas REST para testar o create e o list.

### gRPC
Utilizando o Evans passo a passo:

Acessar o serviço: `evans -r repl --host 127.0.0.1 --port 50051`

Acessar o package: `package pb`

Acessar o service: `service OrderService`

Chamar o create: `call CreateOrder`

Chamar o list: `call ListOrders`

### GraphQL

Acessar pelo navegador: `http://localhost:8080/`

Digitar a mutation e a query e rodar `createOrder` e `queryOrder`

```graphql
mutation createOrder {
  createOrder(input: {id: "xxxxx", Price: 10.2, Tax: 2.0}) {
    id
    Price
    Tax
    FinalPrice
  }
}

query queryOrder {
  order {
    id
    Price
    Tax
    FinalPrice
  }
}
