# Objetivo

Desenvolver um sistema em Go que receba um CEP, identifique a cidade e retorne o clima atual (temperatura em graus Celsius, Fahrenheit e Kelvin). Esse sistema deverá ser publicado no Google Cloud Run.

# Requisitos

- O sistema deve receber um CEP válido de 8 dígitos.
- O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização.
- A partir da localização, o sistema deverá retornar as temperaturas formatadas em:
  - Celsius
  - Fahrenheit
  - Kelvin

## O sistema deve responder adequadamente nos seguintes cenários:

### Em caso de sucesso:
- **Código HTTP**: 200
- **Response Body**:
  ```json
  {
    "temp_C": 28.5,
    "temp_F": 83.3,
    "temp_K": 301.5
  }

### Em caso de falha

- **Se o CEP não for válido (com formato correto)**:
  - **Código HTTP**: 422
  - **Mensagem**: `invalid zipcode`

- **Se o CEP não for encontrado**:
  - **Código HTTP**: 404
  - **Mensagem**: `can not find zipcode`

O deploy deverá ser realizado no Google Cloud Run.

### Dicas

- Utilize a API [viaCEP](https://viacep.com.br/) (ou similar) para encontrar a localização.
- Utilize a API [WeatherAPI](https://www.weatherapi.com/) (ou similar) para consultar as temperaturas desejadas.
- Para realizar a conversão de Celsius para Fahrenheit, utilize a seguinte fórmula:
  - `F = C * 1,8 + 32`
- Para realizar a conversão de Celsius para Kelvin, utilize a seguinte fórmula:
  - `K = C + 273`

  Onde:
  - **F** = Fahrenheit
  - **C** = Celsius
  - **K** = Kelvin

### Entrega

- O código-fonte completo da implementação.
- Testes automatizados demonstrando o funcionamento.
- Utilize Docker/Docker Compose para os testes da aplicação.
- Deploy realizado no Google Cloud Run (free tier) com o endereço ativo para acesso.
