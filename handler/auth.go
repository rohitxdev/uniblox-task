package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rohitxdev/go-api-starter/auth"
	"github.com/rohitxdev/go-api-starter/repo"
)

func (h *Handler) LogOut(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	if sess.IsNew {
		return echo.NewHTTPError(http.StatusBadRequest, "User is not logged in")
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response{Message: "Logged out successfully"})
}

type logInRequest struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required"`
}

func (h *Handler) LogIn(c echo.Context) error {
	req := new(logInRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}
	userEmail := sanitizeEmail(req.Email)
	user, err := h.Repo.GetUserByEmail(c.Request().Context(), userEmail)
	if err != nil {
		if errors.Is(err, repo.ErrUserNotFound) {
			return c.JSON(http.StatusUnauthorized, response{Message: "Invalid credentials"})
		}
		return err
	}
	if !auth.VerifyPassword(req.Password, user.PasswordHash) {
		return c.JSON(http.StatusUnauthorized, response{Message: "Invalid credentials"})
	}
	if _, err = CreateSession(c, h.Config.SessionDuration, user.ID, !h.Config.IsDev); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, response{Message: "Logged in successfully"})
}

type signUpRequest struct {
	Email    string `form:"email" json:"email" validate:"required,email"`
	Password string `form:"password" json:"password" validate:"required"`
}

func (h *Handler) SignUp(c echo.Context) error {
	req := new(signUpRequest)
	if err := bindAndValidate(c, req); err != nil {
		return err
	}
	userEmail := sanitizeEmail(req.Email)
	var userID int
	_, err := h.Repo.GetUserByEmail(c.Request().Context(), userEmail)
	if err == nil {
		return c.JSON(http.StatusBadRequest, response{Message: "User already exists"})
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("Failed to hash password: %w", err)
	}

	if userID, err = h.Repo.CreateUser(c.Request().Context(), userEmail, passwordHash); err != nil {
		return fmt.Errorf("Failed to set password hash: %w", err)
	}
	if err = h.Repo.SetIsVerified(c.Request().Context(), userID, true); err != nil {
		return err
	}
	if _, err = h.Repo.CreateCoupon(c.Request().Context(), userID, "UNIBLOX10", 10); err != nil {
		return err
	}
	if _, err = CreateSession(c, h.Config.SessionDuration, userID, !h.Config.IsDev); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, response{Message: "Signed up successfully"})
}
