# Desafio Cloud Run - API de Clima

Este projeto implementa uma API REST que consulta informações de clima baseadas em CEP (Código de Endereçamento Postal) brasileiro.

## Funcionalidades

- Consulta de temperatura por CEP
- Validação de formato de CEP
- Integração com APIs externas (ViaCEP e WeatherAPI)
- Retorno de temperaturas em Celsius, Fahrenheit e Kelvin

## Estrutura do Projeto

```
desafio-cloudRun/
├── main.go              # Aplicação principal
├── main_test.go         # Testes unitários com testify
├── go.mod               # Dependências do projeto
├── go.sum               # Checksums das dependências
├── Dockerfile           # Configuração para containerização
├── docker-compose.yml   # Configuração para desenvolvimento local
└── deploy-instructions.md # Instruções de deploy
```

## Como Executar

### Localmente

```bash
# Instalar dependências
go mod tidy

# Executar a aplicação
go run main.go
```

A aplicação estará disponível em `http://localhost:8080`

### Com Docker

```bash
# Construir a imagem
docker build -t desafio-cloudrun .

# Executar o container
docker run -p 8080:8080 desafio-cloudrun
```

## Testes

O projeto utiliza o framework [testify](https://github.com/stretchr/testify) para testes unitários, que oferece assertivas mais expressivas e legíveis.

### Executar Testes

```bash
# Executar todos os testes
go test

# Executar testes com detalhes
go test -v

# Executar testes com cobertura
go test -v -cover

# Executar testes com relatório de cobertura detalhado
go test -v -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Cobertura de Testes

Atualmente a cobertura de testes é de **68.0%** das declarações do código.

### Casos de Teste Implementados

1. **TestHandleWeather_Success** - Testa o cenário de sucesso com CEP válido
2. **TestHandleWeather_InvalidCEP** - Testa CEP com formato inválido
3. **TestHandleWeather_NotFound** - Testa CEP não encontrado
4. **TestHandleWeather_InvalidCEPMenor8Digitos** - Testa CEP com menos de 8 dígitos
5. **TestHandleWeather_EmptyCEP** - Testa requisição sem CEP
6. **TestHandleWeather_CEPMaior8Digitos** - Testa CEP com mais de 8 dígitos

### Funcionalidades do Testify Utilizadas

- **assert.Equal()** - Verifica igualdade entre valores
- **assert.InDelta()** - Verifica igualdade com tolerância para números de ponto flutuante
- **assert.Contains()** - Verifica se uma string contém outra
- **require.NoError()** - Verifica que não há erro (falha o teste se houver erro)

## API Endpoints

### GET /weather

Consulta informações de clima para um CEP específico.

**Parâmetros:**

- `cep` (query string): CEP brasileiro (8 dígitos)

**Exemplo de uso:**

```bash
curl "http://localhost:8080/weather?cep=01001000"
```

**Resposta de sucesso:**

```json
{
  "temp_c": 22.5,
  "temp_f": 72.5,
  "temp_k": 295.65
}
```

**Códigos de erro:**

- `422` - CEP inválido
- `404` - CEP não encontrado ou erro na consulta de clima

## Deploy no Google Cloud Run no Free Tier

Este projeto inclui configurações para executar testes automatizados gratuitamente usando GitHub Actions e Google Cloud Run.

### Aplicação Deployada

A aplicação está disponível em produção: [https://weather-service-877639696034.us-central1.run.app/weather?cep=13175658](https://weather-service-877639696034.us-central1.run.app/weather?cep=13175658)

**Exemplo de resposta:**

```json
{ "temp_c": 21.5, "temp_f": 70.7, "temp_k": 294.64 }
```
