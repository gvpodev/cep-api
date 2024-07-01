package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetCEPResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	OpenCEP      string `json:"open-cep"`
}

const (
	url = "https://brasilapi.com.br/api/cep/v1/%s"
)

type BrasilAPIGateway struct {
	Client http.Client
}

func NewBrasilAPIGateway(client http.Client) *BrasilAPIGateway {
	return &BrasilAPIGateway{Client: client}
}

func (gtw *BrasilAPIGateway) GetCEP(cep string) (*GetCEPResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(url, cep), nil)
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

	var response *GetCEPResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
