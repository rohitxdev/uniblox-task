package handler

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func setUpRoutes(e *echo.Echo, svc *Services) {
	h := &Handler{svc}

	e.GET("/metrics", echoprometheus.NewHandler())
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler())
	e.GET("/config", h.GetConfig)
	e.GET("/me", h.GetMe, h.require(RoleUser))
	e.GET("/", h.GetHome)

	auth := e.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/log-in", h.LogIn)
		auth.GET("/log-out", h.LogOut)
	}

	products := e.Group("/products")
	{
		products.GET("", h.GetProducts)
	}

	cart := e.Group("/carts")
	{
		cart.GET("", h.GetCart, h.require(RoleUser))
		cart.POST("/:productId", h.AddToCart, h.require(RoleUser))
		cart.PUT("/:productId/:quantity", h.UpdateCartItemQuantity, h.require(RoleUser))
		cart.DELETE("/:productId", h.DeleteCartItem, h.require(RoleUser))
	}

	orders := e.Group("/orders")
	{
		orders.POST("", h.CreateOrder, h.require(RoleUser))
	}

	coupons := e.Group("/coupons")
	{
		coupons.GET("", h.GetAvailableCoupons, h.require(RoleUser))
		coupons.POST("", h.CreateCoupon, h.require(RoleAdmin))
	}

	admin := e.Group("/_", h.require(RoleAdmin))
	{
		admin.GET("", h.GetAdmin)
		admin.GET("/orders", h.GetAllOrders)
		admin.GET("/coupons", h.GetAllCoupons)
	}
}
