package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rohitxdev/go-api-starter/repo"
)

type GetCartResponse struct {
	Data []repo.CartItem `json:"data"`
}

func (h *Handler) GetCart(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	cart, err := h.Repo.GetCart(c.Request().Context(), user.Id)
	if err != nil {
		if err == repo.ErrCartNotFound {
			return c.JSON(http.StatusNotFound, response{Message: "Cart not found."})
		}
		return err
	}

	return c.JSON(http.StatusOK, GetCartResponse{Data: cart})
}

type AddToCartRequest struct {
	ProductID int `param:"productId" validate:"required"`
	Quantity  int `param:"quantity" validate:"required,min=1"`
}

func (h *Handler) AddToCart(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	var req AddToCartRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	if err := h.Repo.AddToCart(c.Request().Context(), user.Id, req.ProductID, req.Quantity); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response{Message: "Cart item added."})
}

type UpdateCartItemQuantityRequest struct {
	ProductID int `param:"productId" validate:"required"`
	Quantity  int `param:"quantity" validate:"required,min=1"`
}

func (h *Handler) UpdateCartItemQuantity(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	var req UpdateCartItemQuantityRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	if err := h.Repo.UpdateCartItemQuantity(c.Request().Context(), user.Id, req.ProductID, req.Quantity); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response{Message: "Cart item quantity updated."})
}
