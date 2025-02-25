package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) HelloWorldHandler(c echo.Context) error {
	resp := &HelloWorldresponse{
		Message:  "Hello World",
		Message2: "Hello World2",
		Message3: "Hello World3",
		// ErrorMessage: "Error",
		// Error: &Errors{
		// 	Error:  "Error",
		// 	Error2: "Error2",
		// },
	}
	return c.JSON(http.StatusOK, resp)
}
