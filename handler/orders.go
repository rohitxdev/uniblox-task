package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rohitxdev/go-api-starter/repo"
)

type CreateOrderRequest struct {
	CouponCode string `query:"couponCode"`
}

type CreateOrderResponse struct {
	Order *repo.Order `json:"order"`
}

// @Summary Create order
// @Description Create order.
// @Router /orders [post]
// @Security ApiKeyAuth
// @Param couponCode query string false "Coupon code"
// @Success 200 {object} CreateOrderResponse
// @Failure 400 {string} string "invalid coupon"
// @Failure 401 {string} string "invalid session"
func (h *Handler) CreateOrder(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	var req CreateOrderRequest
	req.CouponCode = c.QueryParam("couponCode")

	var coupon *repo.Coupon
	if req.CouponCode != "" {
		coupons, err := h.Repo.GetAvailableCoupons(c.Request().Context(), user.Id)
		if err != nil {
			return err
		}
		for _, c := range coupons {
			if c.Code == req.CouponCode {
				coupon = &c
				break
			}
		}
		if coupon.IsUsed || coupon.UserID != user.Id {
			return c.JSON(http.StatusBadRequest, response{Message: "Invalid coupon"})
		}
	}

	var preTotalAmount int

	cartItems, products, err := h.Repo.GetCartWithProducts(c.Request().Context(), user.Id)
	if err != nil {
		return err
	}

	if len(cartItems) != len(products) {
		return errors.New("cart and products length mismatch")
	}

	for i, cartItem := range cartItems {
		preTotalAmount += cartItem.Quantity * products[i].Price
	}

	order, err := h.Repo.CreateOrder(c.Request().Context(), cartItems, user.Id, preTotalAmount, coupon)
	if err != nil {
		return err
	}

	if err = h.Repo.DiscardCart(c.Request().Context(), user.Id); err != nil {
		return err
	}

	// Create coupon for next order if user is eligible
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		coupons, err := h.Repo.GetAvailableCoupons(ctx, user.Id)
		if err != nil {
			h.Logger.Err(err).Msg("Failed to get available coupons")
			return
		}
		// Return if user already has a coupon
		if len(coupons) > 0 {
			return
		}
		ordersCount, err := h.Repo.GetOrdersCountForUser(ctx, user.Id)
		if err != nil {
			h.Logger.Err(err).Msg("Failed to get orders count for user")
			return
		}

		// Create a coupon if orders count is divisible by 5
		if (ordersCount+1)%5 != 0 {
			return
		}

		if _, err = h.Repo.CreateCoupon(ctx, user.Id, ulid.Make().String(), 10); err != nil {
			h.Logger.Err(err).Msg("Failed to create coupon")
			return
		}
	}()

	return c.JSON(http.StatusCreated, CreateOrderResponse{Order: order})
}

type GetAllOrdersRequest struct {
	Page     int `query:"page"`
	PageSize int `query:"pageSize"`
}

type GetAllOrdersResponse struct {
	Orders []repo.Order `json:"orders"`
}

// @Summary Get all orders
// @Description Get all orders.
// @Router /orders/all [get]
// @Security ApiKeyAuth
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} GetAllOrdersResponse
// @Failure 401 {string} string "invalid session"
func (h *Handler) GetAllOrders(c echo.Context) error {
	var req GetAllOrdersRequest
	var err error
	if err = bindAndValidate(c, &req); err != nil {
		return err
	}

	orders, err := h.Repo.GetAllOrders(c.Request().Context(), req.Page, req.PageSize)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GetAllOrdersResponse{Orders: orders})
}
