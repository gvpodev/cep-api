package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetViaCEPResponse struct {
	Cep          string `json:"cep"`
	Street       string `json:"logradouro"`
	Number       string `json:"number"`
	Unit         string `json:"unit"`
	Neighborhood string `json:"bairro"`
	City         string `json:"localidade"`
	State        string `json:"uf"`
	IBGE         string `json:"ibge"`
	GIA          string `json:"gia"`
	DDD          string `json:"ddd"`
	Siafi        string `json:"siafi"`
}

const viaCEPUrl = "https://viacep.com.br/ws/%s/json/"

type ViaCEPAPIGateway struct {
	client *http.Client
}

func NewViaCEPAPIGateway(c *http.Client) *ViaCEPAPIGateway {
	return &ViaCEPAPIGateway{client: c}
}

func (gtw *ViaCEPAPIGateway) GetViaCEP(cep string) (*GetViaCEPResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(viaCEPUrl, cep), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response *GetViaCEPResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
