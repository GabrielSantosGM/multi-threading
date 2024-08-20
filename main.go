package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BrasilAPI struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	var brasilAPI BrasilAPI
	var viaCEP ViaCep

	channel1 := make(chan BrasilAPI)
	channel2 := make(chan ViaCep)

	zipCode := "09181040"

	go func() {
		requestToBrasilApi(&brasilAPI, zipCode)
		channel1 <- brasilAPI
	}()

	go func() {
		requestToViaCep(&viaCEP, zipCode)
		channel2 <- viaCEP
	}()

	select {
	case brazilAPI := <-channel1:
		fmt.Printf("O retorno foi da Brasil API: %s", brazilAPI)

	case viaCEP := <-channel2:
		fmt.Printf("O retorno foi da ViaCEP: %s", viaCEP)

	case <-time.After(time.Second * 1):
		println("Timeout!")
	}
}

func requestToBrasilApi(brasilAPI *BrasilAPI, zipCode string) error {
	resp, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + zipCode)
	if err != nil {
		fmt.Println(err)
		return err
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := json.Unmarshal(res, &brasilAPI); err != nil {
		fmt.Println(err)
	}
	return nil
}

func requestToViaCep(viaCep *ViaCep, zipCode string) error {
	resp, err := http.Get("https://viacep.com.br/ws/" + zipCode + "/json/")
	if err != nil {
		fmt.Println(err)
		return err
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := json.Unmarshal(res, &viaCep); err != nil {
		fmt.Println(err)
	}
	return nil
}
