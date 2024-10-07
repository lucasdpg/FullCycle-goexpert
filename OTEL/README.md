## Objetivo

Desenvolver um sistema em Go que receba um CEP, identifique a cidade e retorne o clima atual (temperatura em graus Celsius, Fahrenheit e Kelvin) juntamente com a cidade. Esse sistema deverá implementar OTEL (Open Telemetry) e Zipkin.

Baseado no cenário conhecido "Sistema de temperatura por CEP", denominado Serviço B, será incluído um novo projeto, denominado Serviço A.

## Requisitos - Serviço A (responsável pelo input)

- O sistema deve receber um input de 8 dígitos via POST, através do seguinte schema:

  ```json
  {
    "cep": "29902555"
  }
  ```

- O sistema deve validar se o input é válido (contém 8 dígitos) e é uma **STRING**.
- Caso seja válido, será encaminhado para o **Serviço B** via HTTP.
- Caso não seja válido, deve retornar:
  - **Código HTTP**: `422`
  - **Mensagem**: `invalid zipcode`

## Requisitos - Serviço B (responsável pela orquestração)

- O sistema deve receber um **CEP válido de 8 dígitos**.
- O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização. A partir disso, deverá retornar as temperaturas e formatá-las em: **Celsius**, **Fahrenheit**, **Kelvin** juntamente com o nome da localização.
- O sistema deve responder adequadamente nos seguintes cenários:
  - **Em caso de sucesso**:
    - **Código HTTP**: `200`
    - **Response Body**:
      ```json
      {
        "city": "São Paulo",
        "temp_C": 28.5,
        "temp_F": 83.3,
        "temp_K": 301.65
      }
      ```
  - **Em caso de falha, caso o CEP não seja válido (com formato correto)**:
    - **Código HTTP**: `422`
    - **Mensagem**: `invalid zipcode`
  - **Em caso de falha, caso o CEP não seja encontrado**:
    - **Código HTTP**: `404`
    - **Mensagem**: `cannot find zipcode`

## OTEL + Zipkin

- Após a implementação dos serviços, adicione a implementação do **OTEL** + **Zipkin**:
  - Implementar tracing distribuído entre **Serviço A** e **Serviço B**.
  - Utilizar **span** para medir o tempo de resposta do serviço de busca de CEP e de busca de temperatura.

## Dicas

- Utilize a API [viaCEP](https://viacep.com.br/) (ou similar) para encontrar a localização que deseja consultar a temperatura.
- Utilize a API [WeatherAPI](https://www.weatherapi.com/) (ou similar) para consultar as temperaturas desejadas.
- Para realizar a conversão de Celsius para Fahrenheit, utilize a seguinte fórmula:  
  **F = C * 1.8 + 32**
- Para realizar a conversão de Celsius para Kelvin, utilize a seguinte fórmula:  
  **K = C + 273**

  Sendo:  
  - **F** = Fahrenheit  
  - **C** = Celsius  
  - **K** = Kelvin

- Para dúvidas sobre a implementação do **OTEL**, você pode [clicar aqui](#).
- Para implementação de spans, você pode [clicar aqui](#).
- Você precisará utilizar um serviço de **collector do OTEL**.
- Para mais informações sobre **Zipkin**, você pode [clicar aqui](#).

## Entrega

- O código-fonte completo da implementação.
- Documentação explicando como rodar o projeto em ambiente dev.
- Utilize **docker/docker-compose** para que possamos realizar os testes de sua aplicação.

# Como rodar o projeto em ambiente dev.

Para rodar o projeto e fazer os testes basta rodar os comandos abaixo para subir as apps do docker compose e fazer os testes. 

1. Fazer o clone

```
git clone git@github.com:lucasdpg/FullCycle-goexpert.git
```

2. docker compose

Acessar a pasta raiz do projeto e rodar o comando do docker compose

```
cd FullCycle-goexpert/OTEL
```

```
docker compose up -d
```

3. O acesso ao Zipkin pelo browser.

```
http://localhost:9411/
```

4. Comandos exemplos para fazer as chamadas no serviço A e validar o projeto.

Status 200
```
curl -X POST localhost:3030/zipcode -H "Content-Type: application/json" -d '{"cep": "13331630"}'
```

Status 404
```
curl -X POST localhost:3030/zipcode -H "Content-Type: application/json" -d '{"cep": "10101010"}'
```

Status 422
```
curl -X POST localhost:3030/zipcode -H "Content-Type: application/json" -d '{"cep": "1010"}'
```
