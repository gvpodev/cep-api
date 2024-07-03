package fiber

import (
	"cep-api/internal/handler"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func Handlers(h *handler.CEPHandler) http.HandlerFunc {
	r := fiber.New()
	r.Get("/health", healthHandler)
	r.Get("/cep/:cepValue", getCEP(h))
	return adaptor.FiberApp(r)
}

func getCEP(h *handler.CEPHandler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		cep := c.Params("cepValue")
		result := h.GetCEP(cep)
		if result.Fail != nil {
			return c.Status(http.StatusBadRequest).SendString(result.Fail.Err.Error())
		}
		return c.Status(http.StatusOK).JSON(result.Data)
	}
}

func healthHandler(c *fiber.Ctx) error {
	return c.SendString("App is healthy")
}
