package handler

import (
	"cep-api/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CEPHandler struct {
}

func NewCEPHandler() *CEPHandler {
	return &CEPHandler{}
}

const (
	brAPIUrl  = "https://brasilapi.com.br/api/cep/v1/%s"
	viaCEPUrl = "https://viacep.com.br/ws/%s/json/"
)

func (h *CEPHandler) GetCEP(cep string) *model.Result {
	brapiCh := make(chan *model.GetBrAPICEPResponse)
	viacepCh := make(chan *model.GetViaCEPResponse)
	errCh := make(chan *model.ErrResult)

	go func() {
		for {
			brapiCh <- doRequest[model.GetBrAPICEPResponse](errCh, brAPIUrl, cep)
		}
	}()

	go func() {
		for {
			viacepCh <- doRequest[model.GetViaCEPResponse](errCh, viaCEPUrl, cep)
		}
	}()

	for {
		select {
		case err := <-errCh:
			fmt.Println("ERROR")
			return &model.Result{
				Data: nil,
				Fail: err,
			}
		case msg := <-brapiCh:
			fmt.Println("BR API")
			return &model.Result{
				Data: &model.CEPResponse{
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
			return &model.Result{
				Data: &model.CEPResponse{
					CEP:          msg.CEP,
					Street:       msg.Street,
					Neighborhood: msg.Neighborhood,
					City:         msg.City,
					State:        msg.State,
				},
				Fail: nil,
			}
		case <-time.After(1 * time.Second):
			return &model.Result{
				Data: nil,
				Fail: &model.ErrResult{
					Err:        http.ErrHandlerTimeout,
					StatusCode: http.StatusGatewayTimeout,
				},
			}
		}
	}
}

func doRequest[T any](done chan *model.ErrResult, url string, cep string) *T {
	req, err := http.NewRequest("GET", fmt.Sprintf(url, cep), nil)
	if err != nil {
		done <- &model.ErrResult{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
		return nil
	}

	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		done <- &model.ErrResult{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
		return nil
	}
	if res.StatusCode != http.StatusOK {
		done <- &model.ErrResult{
			Err:        errors.New(fmt.Sprintf("Client error: %d", res.StatusCode)),
			StatusCode: res.StatusCode,
		}
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		done <- &model.ErrResult{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
		return nil
	}

	var response *T
	err = json.Unmarshal(body, &response)
	if err != nil {
		done <- &model.ErrResult{
			Err:        err,
			StatusCode: http.StatusInternalServerError,
		}
		return nil
	}

	fmt.Printf("%+v\n", response)

	return response
}
