package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {

	cep := "01153000"

	url1 := "https://viacep.com.br/ws/"
	compl1 := "/json/"
	resultChan1 := make(chan string)

	url2 := "https://brasilapi.com.br/api/cep/v1/"
	compl2 := ""
	resultChan2 := make(chan string)

	go func() {
		body, err := reqCep(cep, url1, compl1)
		if err != nil {
			panic(err)
		}
		// Validar timeout
		//time.Sleep(time.Second * 2)
		resultChan1 <- body
	}()

	go func() {
		body, err := reqCep(cep, url2, compl2)
		if err != nil {
			panic(err)
		}
		// Validar timeout
		// time.Sleep(time.Second * 2)
		resultChan2 <- body
	}()

	select {
	case res1 := <-resultChan1:
		fmt.Println("--- Resposta API viacep ---")
		fmt.Println(string(res1))
	case res2 := <-resultChan2:
		fmt.Println("--- Resposta API brasilapi ---")
		fmt.Println(string(res2))
	case <-time.After(time.Second * 1):
		fmt.Println("--- Timeout ---")
	}

}

func reqCep(cep, url, compl string) (string, error) {

	req, err := http.Get(url + cep + compl)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	if req.StatusCode != http.StatusOK {
		return "", err
	}

	var jsonData map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&jsonData); err != nil {
		return "", err
	}

	response, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(response), nil

}
