package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rohitxdev/go-api-starter/repo"
)

type GetAllCouponsRequest struct {
	Page     string `query:"page"`
	PageSize string `query:"pageSize"`
}

type GetAllCouponsResponse struct {
	Coupons []repo.Coupon `json:"coupons"`
}

// @Summary Get all coupons
// @Description Get all coupons.
// @Router /coupons/all [get]
// @Security ApiKeyAuth
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} GetAllCouponsResponse
// @Failure 401 {string} string "invalid session"
func (h *Handler) GetAllCoupons(c echo.Context) error {
	var req GetAllCouponsRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	page, err := strconv.Atoi(req.Page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response{Message: "Invalid page"})
	}
	pageSize, err := strconv.Atoi(req.PageSize)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response{Message: "Invalid page size"})
	}

	if pageSize <= 0 {
		pageSize = 20
	}
	if page < 0 {
		page = 1
	}

	coupons, err := h.Repo.GetAllCoupons(c.Request().Context(), page, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, GetAllCouponsResponse{Coupons: coupons})
}

type GetAvailableCouponsResponse struct {
	Coupons []repo.Coupon `json:"coupons"`
}

// @Summary Get available coupons
// @Description Get available coupons.
// @Router /coupons [get]
// @Security ApiKeyAuth
// @Success 200 {object} GetAvailableCouponsResponse
// @Failure 401 {string} string "invalid session"
func (h *Handler) GetAvailableCoupons(c echo.Context) error {
	user := getUser(c)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	coupons, err := h.Repo.GetAvailableCoupons(c.Request().Context(), user.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GetAvailableCouponsResponse{Coupons: coupons})
}

type CreateCouponRequest struct {
	DiscountPercent int `json:"discountPercent" validate:"required,min=0,max=100"`
	UserID          int `json:"userId" validate:"required"`
}

type CreateCouponResponse struct {
	Coupon *repo.Coupon `json:"coupon"`
}

// @Summary Create coupon
// @Description Create coupon.
// @Router /coupons [post]
// @Security ApiKeyAuth
// @Param discountPercent path int true "Discount percent"
// @Param userId path int true "User ID"
// @Success 200 {object} CreateCouponResponse
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "invalid session"
func (h *Handler) CreateCoupon(c echo.Context) error {
	var req CreateCouponRequest
	if err := bindAndValidate(c, &req); err != nil {
		return err
	}

	user := getUser(c)
	if user == nil {
		return c.NoContent(http.StatusUnauthorized)
	}

	coupon, err := h.Repo.CreateCoupon(c.Request().Context(), user.ID, ulid.Make().String(), req.DiscountPercent)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, CreateCouponResponse{Coupon: coupon})
}
