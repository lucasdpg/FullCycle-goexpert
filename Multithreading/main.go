package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	cep := "01153000"
	//url1 := "https://viacep.com.br/ws/"
	//compl1 := "/json/"
	url2 := "https://brasilapi.com.br/api/cep/v1/"
	compl2 := ""

	body, err := reqCep(cep, url2, compl2)
	if err != nil {
		fmt.Println("Erro:", err)
		return
	}

	jsonData, err := readBody(body)
	if err != nil {
		fmt.Println("Error to read response:", err)
		return
	}

	response, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		fmt.Println("Error to marshal JSON:", err)
		return
	}

	fmt.Println(string(response))
}

func reqCep(cep, url, compl string) (io.ReadCloser, error) {
	req, err := http.Get(url + cep + compl)
	if err != nil {
		return nil, fmt.Errorf("error to get CEP %v: %v", cep, err)
	}

	if req.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error to get CEP %v: status %v", cep, req.Status)
	}

	return req.Body, nil
}

func readBody(body io.ReadCloser) (map[string]interface{}, error) {

	defer body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error to decode body: %v", err)
	}

	return result, nil

}
