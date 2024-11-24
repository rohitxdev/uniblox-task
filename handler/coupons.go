package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rohitxdev/go-api-starter/repo"
)

type GetAllCouponsRequest struct {
	Page     int `query:"page"`
	PageSize int `query:"pageSize"`
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

	coupons, err := h.Repo.GetAllCoupons(c.Request().Context(), req.Page, req.PageSize)
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

	coupons, err := h.Repo.GetAvailableCoupons(c.Request().Context(), user.Id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, GetAvailableCouponsResponse{Coupons: coupons})
}
