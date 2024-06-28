package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	cep := "01153000"

	url1 := "https://viacep.com.br/ws/"
	compl1 := "/json/"
	resultChan1 := make(chan string)

	//url2 := "https://brasilapi.com.br/api/cep/v1/"
	//compl2 := ""

	go func() {
		body, err := reqCep(cep, url1, compl1)
		if err != nil {
			panic(err)
		}
		resultChan1 <- body
	}()

	res := <-resultChan1

	fmt.Println(string(res))
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
