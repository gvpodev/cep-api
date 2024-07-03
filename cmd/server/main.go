package main

import (
	"cep-api/internal/api"
	"cep-api/internal/handler"
	"cep-api/internal/http/gin"
	"log"
)

func main() {
	h := handler.NewCEPHandler()
	g := gin.Handlers(h)

	err := api.Start("8080", g)
	if err != nil {
		log.Fatalf("error running api: %s", err)
	}
}
