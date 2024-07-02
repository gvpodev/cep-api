package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GetBrAPICEPResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	OpenCEP      string `json:"open-cep"`
}

const brAPIUrl = "https://brasilapi.com.br/api/cep/v1/%s"

type BrasilAPIGateway struct {
	client *http.Client
}

func NewBrasilAPIGateway(c *http.Client) *BrasilAPIGateway {
	return &BrasilAPIGateway{client: c}
}

func (gtw *BrasilAPIGateway) GetBrAPI(cep string) (*GetBrAPICEPResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(brAPIUrl, cep), nil)
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

	var response *GetBrAPICEPResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
