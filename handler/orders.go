package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rohitxdev/go-api-starter/repo"
)

func (h *Handler) GetOrders(c echo.Context) error {
	return nil
}

type CreateOrderResponse struct {
	Order *repo.Order `json:"order"`
}

func (h *Handler) CreateOrder(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	var totalAmount int

	cartItems, products, err := h.Repo.GetCartWithProducts(c.Request().Context(), user.Id)
	if err != nil {
		return err
	}

	if len(cartItems) != len(products) {
		return errors.New("cart and products length mismatch")
	}

	for i, cartItem := range cartItems {
		totalAmount += cartItem.Quantity * products[i].Price
	}

	order, err := h.Repo.CreateOrder(c.Request().Context(), cartItems, user.Id, totalAmount)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, CreateOrderResponse{Order: order})
}
