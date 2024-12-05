package handler

import (
	"fmt"
	"net/http"

	"github.com/Ayobami0/chatter_box_server/src/errors"
	"github.com/Ayobami0/chatter_box_server/src/model"
	"github.com/Ayobami0/chatter_box_server/src/service"
	"github.com/Ayobami0/chatter_box_server/src/utils"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	s service.AuthService
}

func NewAuthHandler(service service.AuthService) AuthHandler {
	return AuthHandler{service}
}

func (a *AuthHandler) UserLogin(c echo.Context) error {
	var login model.UserLogin
	if err := c.Bind(&login); err != nil {
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}

	u, err := a.s.VerifyUser(login.Password, login.EmailOrUsername)

	if err != nil {
		switch err.(type) {
		case errors.ErrPasswordIncorrect:
			return utils.ErrorJson(c, http.StatusBadRequest, err.Error())
		case errors.ErrUserNotExist:
			return utils.ErrorJson(c, http.StatusNotFound, err.Error())
		}
	}

	token, err := a.s.GenenerateJWTToken(u.ID, u.Username, u.Email, u.ProfileUrl)

	if err != nil {
		return utils.ErrorJson(c, http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, fmt.Sprintf("Bearer %s", token))
}

func (a *AuthHandler) UserLogout(c echo.Context) error {
	return c.String(http.StatusPermanentRedirect, "LOGOUT JWT SESSION STRING")
}

func (a *AuthHandler) UserCreate(c echo.Context) error {
	var create model.UserCreate

	if err := c.Bind(&create); err != nil {
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}

	user, err := a.s.RegisterUser(create)

	if err != nil {
		switch err.(type) {
		case errors.ErrUserExist:
			return utils.ErrorJson(c, http.StatusBadRequest, err.Error())
		}
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}

	token, err := a.s.GenenerateJWTToken(user.ID, user.Username, user.Email, user.ProfileUrl)

	if err != nil {
		return utils.ErrorJson(c, http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, fmt.Sprintf("Bearer %s", token))
}
