package http

import (
	"net/http"

	"dangquang9a/go-location/internal/auth"
	"dangquang9a/go-location/internal/auth/presenter"
	"dangquang9a/go-location/utils"

	"github.com/labstack/echo/v4"
)

type authHandler struct {
	useCase auth.UseCase
}

func NewAuthHandler(useCase auth.UseCase) auth.Handler {
	return &authHandler{
		useCase: useCase,
	}
}

func (h *authHandler) SignUp() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &presenter.SignUpInput{}
		if err := utils.ReadRequest(c, input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}

		user, err := h.useCase.SignUp(c.Request().Context(), input.Username, input.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusCreated, presenter.SignUpResponse{Id: user.Id, Username: user.Username})
	}
}

func (h *authHandler) SignIn() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &presenter.LoginInput{}
		if err := utils.ReadRequest(c, input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest)
		}
		token, err := h.useCase.SignIn(c.Request().Context(), input.Username, input.Password)
		if err != nil {
			if err == auth.ErrUserNotFound {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			if err == auth.ErrWrongPassword {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.JSON(http.StatusOK, presenter.LogInResponse{Token: token})
	}
}
