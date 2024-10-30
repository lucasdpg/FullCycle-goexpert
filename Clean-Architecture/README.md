# Clean Architecture - Order Management System

## Desafio Proposto

Bem-vindo, Devs!

Chegou a hora de colocar em prática seus conhecimentos. O desafio consiste em criar o caso de uso para a listagem das ordens (`Orders`). Essa listagem deve ser implementada utilizando:

- **Endpoint REST:** `GET /order`
- **Serviço gRPC:** `ListOrders`
- **Query GraphQL:** `ListOrders`

Além disso, não se esqueça de criar as migrações necessárias e o arquivo `api.http` contendo as requisições para criar e listar as ordens.

Para a criação do banco de dados, utilize o Docker. O comando `docker compose up` deverá configurar e iniciar tudo automaticamente, incluindo o banco de dados.

Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.


# Como rodar o projeto e os testes

1. **Clone o repositório e navegue até o diretório do projeto:**

    ```bash
    git clone git@github.com:lucasdpg/FullCycle-goexpert.git && cd FullCycle-goexpert/Clean-Architecture
    ```

2. **Inicie o MySQL o RabbitMQ e Rodar a Migration:**

    ```bash
    docker-compose up -d
    ```

    > **Nota:** Pode ser necessário aguardar alguns minutos até que o banco de dados esteja pronto para o próximo passo.

    > **Nota:** Em alguns casos, a conexão com o banco de dados pode falhar. Caso isso ocorra, reinicie o banco de dados com:

    ```bash
    docker-compose down && docker-compose up -d
    ```

3. **Inicie o Projeto:**

    ```bash
    cd cmd/ordersystem && go run main.go wire_gen.go
    ```

# Validação do Desafio

Se tudo foi configurado corretamente, os serviços estarão disponíveis nas seguintes portas:

- **Web:** `localhost:8000`
- **gRPC:** `localhost:50051`
- **GraphQL:** `localhost:8080`

### Testando a Aplicação

#### WebServer

As requisições REST para criar e listar ordens podem ser testadas utilizando os arquivos `api/create_order.http` e `api/list_order.http`.

#### gRPC

Para testar via gRPC, utilize o [Evans](https://github.com/ktr0731/evans):

1. Conecte-se ao serviço:

    ```bash
    evans -r repl --host 127.0.0.1 --port 50051
    ```

2. Navegue até o package:

    ```bash
    package pb
    ```

3. Acesse o serviço:

    ```bash
    service OrderService
    ```

4. Chame as funções:

    - **Criar Ordem:** `call CreateOrder`
    - **Listar Ordens:** `call ListOrders`

#### GraphQL

Para acessar o serviço GraphQL:

1. Abra o navegador e acesse: `http://localhost:8080/`
2. Digite a mutation e a query para rodar as operações `createOrder` e `queryOrder`:

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
    ```
