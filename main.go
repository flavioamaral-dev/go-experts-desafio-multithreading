package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	fmt.Println("Digite o CEP: ")
	reader := bufio.NewReader(os.Stdin)
	cep, errorOnRead := reader.ReadString('\n')

	if errorOnRead != nil {
		print("Entre com uma cep válido!\n")
		return
	}
	cepReadyFromSearch := formatarCep(cep)

	_, errConvertToInt := strconv.Atoi(cepReadyFromSearch)
	if errConvertToInt != nil {
		print("Entre com uma cep válido!\n")
		return
	}

	realizarConsulta(cepReadyFromSearch)
}

func realizarConsulta(cep string) {
	viacepChannel := make(chan string)
	brasilapiChannel := make(chan string)

	go retornarDadosApi("https://brasilapi.com.br/api/cep/v1/"+cep, viacepChannel)
	go retornarDadosApi("http://viacep.com.br/ws/"+cep+"/json/", brasilapiChannel)

	select {
	case msg1 := <-viacepChannel:
		print("\nViaCep => " + msg1 + "\n")
	case msg2 := <-brasilapiChannel:
		print("\nBrasilAPI => " + msg2 + "\n")
	case <-time.After(time.Second):
		print("TimeOut\n")
	}
}

func retornarDadosApi(url string, c chan<- string) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("Não foi possível criar a requisição: %s\n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Erro ao fazer a requisição: %s\n", err)
	}

	resBody, err := io.ReadAll(res.Body)

	if err == nil {
		c <- string(resBody)
	}

}

func formatarCep(cep string) string {
	aux := strings.TrimSpace(cep)
	aux = strings.Replace(aux, ".", "", -1)
	aux = strings.Replace(aux, "-", "", -1)

	return aux
}
