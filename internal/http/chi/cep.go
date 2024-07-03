package chi

import (
	"cep-api/internal/handler"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func Handlers(h *handler.CEPHandler) http.Handler {
	r := chi.NewRouter()
	r.Get("/health", healthHandler)
	r.Get("/cep/{cepValue}", func(w http.ResponseWriter, r *http.Request) {
		getCEP(h, w, r)
	})

	return r
}

func getCEP(h *handler.CEPHandler, w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cepValue")
	result := h.GetCEP(cep)
	if result.Fail != nil {
		http.Error(w, result.Fail.Err.Error(), result.Fail.StatusCode)
		return
	}

	err := json.NewEncoder(w).Encode(result.Data)
	if err != nil {
		http.Error(w, "Internal error: Encoding.", http.StatusInternalServerError)
		return
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("App is healthy"))
	w.WriteHeader(http.StatusOK)
}
