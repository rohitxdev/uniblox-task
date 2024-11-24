package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rohitxdev/go-api-starter/repo"
)

type GetProductsResponse struct {
	Products []repo.Product `json:"products"`
}

// @Summary Get products
// @Description Get products.
// @Router /products [get]
// @Security ApiKeyAuth
// @Success 200 {object} GetProductsResponse
func (h *Handler) GetProducts(c echo.Context) error {
	products, err := h.Repo.GetProducts(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GetProductsResponse{Products: products})
}
