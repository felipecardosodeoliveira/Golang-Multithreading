package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Address represents the structure of an address.
type Address struct {
	Address interface{} // The address information.
	Origem  string      // The origin of the address data.
}

// searchAddress retrieves address information from the given API URL and sends it to the channel.
func searchAddress(apiURL string, ch chan<- Address) {
	resp, err := http.Get(apiURL)
	if err != nil {
		ch <- Address{nil, fmt.Sprintf("Erro ao fazer a requisição para %s: %s", apiURL, err)}
		return
	}
	defer resp.Body.Close()

	var addressSearch interface{}
	if err := json.NewDecoder(resp.Body).Decode(&addressSearch); err != nil {
		ch <- Address{nil, fmt.Sprintf("Erro ao decodificar a resposta de %s: %s", apiURL, err)}
		return
	}

	ch <- Address{addressSearch, apiURL}
}

// printPrettyJSON prints the JSON data in a pretty format.
func printPrettyJSON(address interface{}) {
	jsonData, err := json.MarshalIndent(address, "", "  ")
	if err != nil {
		fmt.Printf("Erro ao codificar o JSON: %s\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

func main() {
	cep := "81490420"
	ch := make(chan Address, 2)

	apiViaCep := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	apiBrasil := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%scep", cep)

	go searchAddress(apiViaCep, ch)
	go searchAddress(apiBrasil, ch)

	select {
	case address := <-ch:
		if address.Address != nil {
			fmt.Println("Endereço:")
			printPrettyJSON(address.Address)
			fmt.Printf("Origem: %s\n", address.Origem)
		} else {
			fmt.Println(address.Origem)
		}
	case <-time.After(1 * time.Second):
		fmt.Println("Erro: Timeout de 1 segundo excedido.")
	}
}
