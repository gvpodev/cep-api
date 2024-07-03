package gin

import (
	"cep-api/internal/handler"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handlers(h *handler.CEPHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/health", healthHandler)
	r.GET("/cep/:cepValue", getCEP(h))

	return r
}

func getCEP(h *handler.CEPHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		cep := c.Param("cepValue")
		result := h.GetCEP(cep)
		if result.Fail != nil {
			c.JSON(result.Fail.StatusCode, result.Fail.Err.Error())
			return
		}
		c.JSON(http.StatusOK, result.Data)
	}
}

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "App is healthy")
}
