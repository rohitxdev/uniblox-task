package handler

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type role string

const (
	RoleUser  role = "user"
	RoleAdmin role = "admin"
)

var roles = map[role]uint8{
	RoleUser:  1,
	RoleAdmin: 2,
}

func (h *Handler) require(r role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("session", c)
			if err != nil {
				return err
			}
			userID, ok := sess.Values["userId"].(int)
			if !ok {
				return echo.ErrUnauthorized
			}
			user, err := h.Repo.GetUserById(c.Request().Context(), userID)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err)
			}
			if user.AccountStatus != "active" {
				return echo.ErrForbidden
			}
			if roles[role(user.Role)] < roles[role(r)] {
				return echo.ErrForbidden
			}
			c.Set("user", user)
			return next(c)
		}
	}
}
