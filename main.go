package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)


const weatherAPIKey = "4a7f90bec7f2460f8d5170734252207"

type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	UF string `json:"uf"`
	Erro bool `json:"erro"`
}

type WeatherResponse struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

var (
	viaCEPBaseURL     = "https://viacep.com.br/ws"
	weatherAPIBaseURL = "https://api.weatherapi.com/v1"
)


func main() {
	http.HandleFunc("/weather", handleWeatherRequest)
	log.Println("Iniciando aplicação em :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWeatherRequest(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	
	if validatorCep, _ := regexp.MatchString(`^\d{8}$`, cep); !validatorCep {
		log.Printf("CEP inválido: %s", cep)
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	viaCEPURL := fmt.Sprintf("%s/%s/json/", viaCEPBaseURL, cep)
	log.Printf("Consultando ViaCEP: %s", viaCEPURL)
	
	viaCEPResp, err := http.Get(viaCEPURL)
	if err != nil {
		log.Printf("Erro ao consultar ViaCEP: %v", err)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	defer viaCEPResp.Body.Close()

	log.Printf("Status code ViaCEP: %d", viaCEPResp.StatusCode)

	body, _ := io.ReadAll(viaCEPResp.Body)
	log.Printf("Resposta ViaCEP: %s", string(body))
	
	var viaCEP ViaCEPResponse
	if err := json.Unmarshal(body, &viaCEP); err != nil {
		log.Printf("Erro ao fazer parse da resposta ViaCEP: %v", err)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	// Verifica se o CEP foi encontrado (ViaCEP retorna erro: true quando não encontra)
	if viaCEP.Erro {
		log.Printf("CEP não encontrado no ViaCEP: %s", cep)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	// Verifica se os campos necessários estão preenchidos
	if viaCEP.Localidade == "" || viaCEP.UF == "" {
		log.Printf("CEP encontrado mas dados incompletos: localidade=%s, uf=%s", viaCEP.Localidade, viaCEP.UF)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	location := fmt.Sprintf("%s,%s", viaCEP.Localidade, viaCEP.UF)
	encodedLocation := url.QueryEscape(location)

	weatherURL:= fmt.Sprintf("%s/current.json?key=%s&q=%s", weatherAPIBaseURL, weatherAPIKey, encodedLocation)
	weatherResp, err := http.Get(weatherURL)
	if err != nil || weatherResp.StatusCode != 200 {
		http.Error(w, "can not find weather", http.StatusNotFound)
		return
	}
	defer weatherResp.Body.Close()

	var weatherData WeatherAPIResponse
	if err := json.NewDecoder(weatherResp.Body).Decode(&weatherData); err != nil {
		http.Error(w, "weather parser error", http.StatusInternalServerError)
		return
	}

	tempC := weatherData.Current.TempC
	tempF := tempC * 1.8 + 32
	tempK := tempC + 273.15

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(WeatherResponse{
		TempC: round(tempC),
		TempF: round(tempF),
		TempK: round(tempK),
	})
}

func round(value float64) float64 {
	return float64(int(value*100)) / 100
}

