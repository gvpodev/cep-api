package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type GetBrAPICEPResponse struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	OpenCEP      string `json:"open-cep"`
}

type GetViaCEPResponse struct {
	CEP          string `json:"cep"`
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

type CEPResponse struct {
	CEP          string
	Street       string
	Neighborhood string
	City         string
	State        string
}

type Result struct {
	Data *CEPResponse
	Fail *ErrResult
}

type ErrResult struct {
	Err        error
	StatusCode int
}

type CEPHandler struct {
}

func NewCEPHandler() *CEPHandler {
	return &CEPHandler{}
}

const (
	brAPIUrl  = "https://brasilapi.com.br/api/cep/v1/%s"
	viaCEPUrl = "https://viacep.com.br/ws/%s/json/"
)

func (h *CEPHandler) GetCEP(cep string) *Result {
	brapiCh := make(chan *GetBrAPICEPResponse)
	viacepCh := make(chan *GetViaCEPResponse)
	errCh := make(chan *ErrResult)

	go func() {
		for {
			brapiCh <- doRequest[GetBrAPICEPResponse](errCh, brAPIUrl, cep)
		}
	}()

	go func() {
		for {
			viacepCh <- doRequest[GetViaCEPResponse](errCh, viaCEPUrl, cep)
		}
	}()

	for {
		select {
		case err := <-errCh:
			fmt.Println("ERROR")
			return &Result{
				Data: nil,
				Fail: err,
			}
		case msg := <-brapiCh:
			fmt.Println("BR API")
			return &Result{
				Data: &CEPResponse{
					CEP:          msg.CEP,
					Street:       msg.Street,
					Neighborhood: msg.Neighborhood,
					City:         msg.City,
					State:        msg.State,
				},
				Fail: nil,
			}
		case msg := <-viacepCh:
			fmt.Println("VIA CEP")
			return &Result{
				Data: &CEPResponse{
					CEP:          msg.CEP,
					Street:       msg.Street,
					Neighborhood: msg.Neighborhood,
					City:         msg.City,
					State:        msg.State,
				},
				Fail: nil,
			}
		case <-time.After(1 * time.Second):
			return &Result{
				Data: nil,
				Fail: &ErrResult{
					Err:        http.ErrHandlerTimeout,
					StatusCode: http.StatusGatewayTimeout,
				},
			}
		}
	}
}

func doRequest[T any](done chan *ErrResult, url string, cep string) *T {
	req, err := http.NewRequest("GET", fmt.Sprintf(url, cep), nil)
	if err != nil {
		done <- &ErrResult{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
		return nil
	}

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		done <- &ErrResult{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
		return nil
	}
	if res.StatusCode != http.StatusOK {
		done <- &ErrResult{
			Err:        errors.New(fmt.Sprintf("Client error: %d", res.StatusCode)),
			StatusCode: res.StatusCode,
		}
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		done <- &ErrResult{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
		return nil
	}

	var response *T
	err = json.Unmarshal(body, &response)
	if err != nil {
		done <- &ErrResult{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
		return nil
	}

	return response
}
