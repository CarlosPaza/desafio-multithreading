package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Cidade      string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilApi struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"street"`
	Bairro     string `json:"neighborhood"`
	Cidade     string `json:"city"`
	Uf         string `json:"state"`
	Service    string `json:"service"`
}

const DefaultCep = "88353541"

func main() {
	channelViaCEP := make(chan ViaCEP)
	channelBrasilApi := make(chan BrasilApi)

	var cep string
	if len(os.Args[1:]) > 0 {
		cep = os.Args[1]
	} else {
		cep = DefaultCep
	}

	go GetViaCEP(cep, channelViaCEP)
	go GetBrasilApi(cep, channelBrasilApi)

	select {
	case response := <-channelViaCEP:
		printResponse("ViaCEP", response)

	case response := <-channelBrasilApi:
		printResponse("BrasilApi", response)

	case <-time.After(time.Second):
		fmt.Printf("TimeOut")

	}
}

func GetViaCEP(cep string, ch chan<- ViaCEP) {
	var viaCEP ViaCEP
	RequestAPI("https://viacep.com.br/ws/"+cep+"/json/", &viaCEP)
	ch <- viaCEP
}

func GetBrasilApi(cep string, ch chan<- BrasilApi) {
	var brasilApi BrasilApi
	RequestAPI("https://brasilapi.com.br/api/cep/v1/"+cep, &brasilApi)
	ch <- brasilApi
}

func RequestAPI(url string, response interface{}) error {
	req, err := http.Get(url)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	}
	return nil
}

func printResponse(api string, response interface{}) {
	fmt.Print("*** Api: " + api + " *** \n")
	value := fmt.Sprintf("%#v", response)
	value = strings.Split(value, "{")[1]
	value = strings.Replace(value, "}", "", -1)
	fmt.Print(strings.Replace(value, ", ", "\n", -1))
}
