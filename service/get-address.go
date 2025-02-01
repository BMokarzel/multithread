package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	NeighborHodd string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCEPResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func GetAddressHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	cep := query.Get("cep")
	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no zip code found"))
	}

	resOne, resTwo, err := GetAddress(cep)
	if err == errors.New("timeout") {
		w.WriteHeader(http.StatusRequestTimeout)
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	if resOne != nil {
		err := json.NewEncoder(w).Encode(resOne)
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	if resTwo != nil {
		err := json.NewEncoder(w).Encode(resTwo)
		if err != nil {
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	w.WriteHeader(http.StatusOK)

}

func GetAddress(cep string) (*BrasilAPIResponse, *ViaCEPResponse, error) {

	brasilAPIChan := make(chan BrasilAPIResponse)

	viaCEPChan := make(chan ViaCEPResponse)

	go BrasilAPI(cep, brasilAPIChan)

	go ViaCEP(cep, viaCEPChan)

	select {
	case brasilRes := <-brasilAPIChan:
		fmt.Printf("Para o CEP %s temos:\nEstado: %s,\nCidade: %s,\nBairro: %s,\nRua: %s,\n\nServiço de consulta utilizado: BrasilAPI com ajuda de %s", brasilRes.Cep, brasilRes.State, brasilRes.City, brasilRes.NeighborHodd, brasilRes.Street, brasilRes.Service)
		return &brasilRes, nil, nil

	case viaCEPRes := <-viaCEPChan:
		fmt.Printf("Para o CEP %s temos:\nLogradouro: %s,\nComplemento: %s,\nUnidade: %s,\nBairro: %s,\nLocalidade: %s,\nUF: %s,\nEstado: %s,\nRegião: %s,\nIBGE: %s,\nGIA: %s,\nDDD: %s,\nSiafi: %s,\n\nServiço de consulta feito por: VIACEP", viaCEPRes.Cep, viaCEPRes.Logradouro, viaCEPRes.Complemento, viaCEPRes.Unidade, viaCEPRes.Bairro, viaCEPRes.Localidade, viaCEPRes.Uf, viaCEPRes.Estado, viaCEPRes.Regiao, viaCEPRes.Ibge, viaCEPRes.Gia, viaCEPRes.Ddd, viaCEPRes.Siafi)
		return nil, &viaCEPRes, nil

	case <-time.After(time.Second * 1):
		fmt.Println("timeout")
		return nil, nil, errors.New("timeout")
	}
}

func BrasilAPI(cep string, c chan BrasilAPIResponse) {

	req, err := http.NewRequest("GET", "https://brasilapi.com.br/api/cep/v1/"+cep, nil)
	if err != nil {
		return
	}

	fmt.Println("dentro de brasil api, pré req")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	var response BrasilAPIResponse

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return
	}

	c <- response
}

func ViaCEP(cep string, c chan ViaCEPResponse) {

	req, err := http.NewRequest("GET", "http://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {
		return
	}

	fmt.Println("dentro de viacep, pré req")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()

	var response ViaCEPResponse

	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return
	}

	c <- response
}
