package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rohitxdev/go-api-starter/repo"
)

type GetCartResponse struct {
	Cart []repo.CartItem `json:"cart"`
}

// @Summary Get cart
// @Description Get cart.
// @Router /carts [get]
// @Security ApiKeyAuth
// @Success 200 {object} GetCartResponse
// @Failure 401 {string} string "invalid session"
func (h *Handler) GetCart(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	cart, err := h.Repo.GetCart(c.Request().Context(), user.ID)
	if err != nil {
		if err == repo.ErrCartNotFound {
			return c.JSON(http.StatusNotFound, response{Message: "Cart not found."})
		}
		return err
	}

	return c.JSON(http.StatusOK, GetCartResponse{Cart: cart})
}

type AddToCartRequest struct {
	ProductID int `param:"productId" validate:"required"`
}

// @Summary Add to cart
// @Description Add to cart.
// @Router /carts/{productId} [post]
// @Security ApiKeyAuth
// @Param productId path int true "Product ID"
// @Success 200 {object} response
// @Failure 400 {string} string "invalid product"
// @Failure 401 {string} string "invalid session"
func (h *Handler) AddToCart(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	var req AddToCartRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	if err := h.Repo.AddToCart(c.Request().Context(), user.ID, req.ProductID); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response{Message: "Cart item added."})
}

type UpdateCartItemQuantityRequest struct {
	ProductID int `param:"productId" validate:"required"`
	Quantity  int `param:"quantity" validate:"required,min=1"`
}

// @Summary Update cart item quantity
// @Description Update cart item quantity.
// @Router /carts/{productId}/{quantity} [put]
// @Security ApiKeyAuth
// @Param productId path int true "Product ID"
// @Param quantity path int true "Quantity"
// @Success 200 {object} response
// @Failure 400 {string} string "invalid product or quantity"
// @Failure 401 {string} string "invalid session"
func (h *Handler) UpdateCartItemQuantity(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	var req UpdateCartItemQuantityRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	if err := h.Repo.UpdateCartItemQuantity(c.Request().Context(), user.ID, req.ProductID, req.Quantity); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response{Message: "Cart item quantity updated."})
}

type DeleteCartItemRequest struct {
	ProductID int `param:"productId" validate:"required"`
}

// @Summary Delete cart item
// @Description Delete cart item.
// @Router /carts/{productId} [delete]
// @Security ApiKeyAuth
// @Param productId path int true "Product ID"
// @Success 200 {object} response
// @Failure 400 {string} string "invalid product"
// @Failure 401 {string} string "invalid session"
func (h *Handler) DeleteCartItem(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	var req DeleteCartItemRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	if err := h.Repo.DeleteCartItem(c.Request().Context(), user.ID, req.ProductID); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response{Message: "Cart item deleted."})
}
