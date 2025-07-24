package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockViaCEPServer cria um servidor de teste mock para simular a API do ViaCEP

func mockViaCEPServer(t *testing.T, response string, status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Write([]byte(response))
	}))
}


 // mockWeatherAPIServer cria um servidor de teste mock para simular a API do WeatherAPI

func mockWeatherAPIServer(t *testing.T, response string, status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		io.WriteString(w, response)
	}))
}


 //TestHandleWeather_Success testa o cenário de sucesso da API de clima

 func TestHandleWeather_Success(t *testing.T) {
	// Mock de resposta do ViaCEP
	viacep := mockViaCEPServer(t, `{"localidade": "São Paulo", "uf": "SP", "erro": false}`, 200)
	defer viacep.Close()

	// Mock de resposta da WeatherAPI
	weather := mockWeatherAPIServer(t, `{"current": {"temp_c": 22.5}}`, 200)
	defer weather.Close()

	// Configuração dos URLs mock
	originalViaCEP := viaCEPBaseURL
	originalWeatherAPI := weatherAPIBaseURL
	viaCEPBaseURL = viacep.URL
	weatherAPIBaseURL = weather.URL
	defer func() {
		viaCEPBaseURL = originalViaCEP
		weatherAPIBaseURL = originalWeatherAPI
	}()

	// Criação da requisição de teste
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=01001000", nil)
	rec := httptest.NewRecorder()

	// Execução da função de teste
	handleWeatherRequest(rec, req)

	// Verificações usando testify
	assert.Equal(t, http.StatusOK, rec.Code, "Status code deve ser 200")
	
	var resp WeatherResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err, "Não deve haver erro ao decodificar a resposta JSON")
	
	assert.Equal(t, 22.5, resp.TempC, "Temperatura em Celsius deve ser 22.5")
	assert.Equal(t, 72.5, resp.TempF, "Temperatura em Fahrenheit deve ser calculada corretamente")
	assert.InDelta(t, 295.65, resp.TempK, 0.01, "Temperatura em Kelvin deve ser calculada corretamente (com tolerância de 0.01)")
}

// TestHandleWeather_InvalidCEP testa o cenário de CEP inválido

func TestHandleWeather_InvalidCEP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=abcd12", nil)
	rec := httptest.NewRecorder()

	handleWeatherRequest(rec, req)

	// Verificações usando testify
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, "Status code deve ser 422 para CEP inválido")
	
	body := strings.TrimSpace(rec.Body.String())
	assert.Equal(t, "invalid zipcode", body, "Mensagem de erro deve ser 'invalid zipcode'")
}

//TestHandleWeather_NotFound testa o cenário de CEP não encontrado

func TestHandleWeather_NotFound(t *testing.T) {
	viacep := mockViaCEPServer(t, `{"erro": true}`, 200)
	defer viacep.Close()

	originalViaCEP := viaCEPBaseURL
	viaCEPBaseURL = viacep.URL 
	defer func() {
		viaCEPBaseURL = originalViaCEP
	}()

	req := httptest.NewRequest(http.MethodGet, "/weather?cep=99999999", nil)
	rec := httptest.NewRecorder()

	handleWeatherRequest(rec, req)

	// Verificações usando testify
	assert.Equal(t, http.StatusNotFound, rec.Code, "Status code deve ser 404 para CEP não encontrado")
	
	body := strings.TrimSpace(rec.Body.String())
	assert.Equal(t, "can not find zipcode", body, "Mensagem de erro deve ser 'can not find zipcode'")
}

//TestHandleWeather_InvalidCEPMenor8Digitos testa o cenário de CEP com menos de 8 dígitos

func TestHandleWeather_InvalidCEPMenor8Digitos(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=010", nil)
	rec := httptest.NewRecorder()

	handleWeatherRequest(rec, req)

	// Verificações usando testify
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, "Status code deve ser 422 para CEP com formato inválido")
	
	body := rec.Body.String()
	assert.Contains(t, body, "invalid zipcode", "Mensagem de erro deve conter 'invalid zipcode'")
}

// TestHandleWeather_EmptyCEP testa o cenário de CEP vazio

func TestHandleWeather_EmptyCEP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	rec := httptest.NewRecorder()

	handleWeatherRequest(rec, req)

	// Verificações usando testify
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, "Status code deve ser 422 para CEP vazio")
	
	body := strings.TrimSpace(rec.Body.String())
	assert.Equal(t, "invalid zipcode", body, "Mensagem de erro deve ser 'invalid zipcode'")
}

// TestHandleWeather_CEPMaior8Digitos testa o cenário de CEP com mais de 8 dígitos

func TestHandleWeather_CEPMaior8Digitos(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/weather?cep=010010001", nil)
	rec := httptest.NewRecorder()

	handleWeatherRequest(rec, req)

	// Verificações usando testify
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code, "Status code deve ser 422 para CEP com mais de 8 dígitos")
	
	body := strings.TrimSpace(rec.Body.String())
	assert.Equal(t, "invalid zipcode", body, "Mensagem de erro deve ser 'invalid zipcode'")
}