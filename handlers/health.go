package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, h.db.Health())
}
