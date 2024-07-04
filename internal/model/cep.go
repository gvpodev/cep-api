package model

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
