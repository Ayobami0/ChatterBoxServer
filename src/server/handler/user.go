package handler

import (
	"net/http"
	"strconv"

	app_err "github.com/Ayobami0/chatter_box_server/src/errors"
	"github.com/Ayobami0/chatter_box_server/src/service"
	"github.com/Ayobami0/chatter_box_server/src/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	s service.UserService
}

func NewUserHandler(service service.UserService) UserHandler {
	return UserHandler{service}
}

func (u *UserHandler) UsersMeGet(c echo.Context) error {
	return nil
}

func (u *UserHandler) UsersGet(c echo.Context) error {
	qParams := c.QueryParams()

	search := qParams.Get("q")
	page, err := strconv.Atoi(qParams.Get("page"))
	if err != nil {
		page = 1
	}
	count, err := strconv.Atoi(qParams.Get("count"))
	if err != nil {
		count = 10
	}

	users, err := u.s.QueryUsersByUsername(search, page, count)

	if err != nil {
		return utils.ErrorJson(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func (u *UserHandler) UserConversationsGet(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(*service.UserAccountClaim)

	conversations, err := u.s.GetUserConversations(claim.UserId)

	if err != nil {
		return utils.ErrorJson(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, conversations)
}

func (u *UserHandler) UserRequestsGet(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(*service.UserAccountClaim)

	req, err := u.s.GetInviteRequests(claim.UserId)

	if err != nil {
		switch err.(type) {
		case app_err.ErrUserNotExist:
			return utils.ErrorJson(c, http.StatusNotFound, err.Error())
    }
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}

	return c.JSON(http.StatusOK, req)
}

func (u *UserHandler) UserRequestsAccept(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(*service.UserAccountClaim)

  reqId := c.Param("id")

	err := u.s.AcceptRequest(claim.UserId, reqId)

	if err != nil {
		switch err.(type) {
		case app_err.ErrUserNotExist:
			return utils.ErrorJson(c, http.StatusNotFound, err.Error())
		case app_err.ErrNoSuchRequest:
			return utils.ErrorJson(c, http.StatusNotFound, err.Error())
    }
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}
	return c.String(http.StatusOK, "OK")
}

func (u *UserHandler) UserRequestsReject(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claim := user.Claims.(*service.UserAccountClaim)

  reqId := c.Param("id")

	err := u.s.DeclineRequest(claim.UserId, reqId)

	if err != nil {
		switch err.(type) {
		case app_err.ErrUserNotExist:
			return utils.ErrorJson(c, http.StatusNotFound, err.Error())
		case app_err.ErrNoSuchRequest:
			return utils.ErrorJson(c, http.StatusNotFound, err.Error())
    }
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}
	return c.String(http.StatusOK, "OK")
}

func (u *UserHandler) UserInvite(c echo.Context) error {
  user := c.Get("user").(*jwt.Token)
  claim := user.Claims.(*service.UserAccountClaim)

  toId := c.Param("id") // Id of user to send the request to

  if toId == claim.UserId {
    return utils.ErrorJson(c, http.StatusBadRequest, "cannot send invite to self")
  }

	err := u.s.SendInviteToUser(claim.UserId, toId)

	if err != nil {
		switch err.(type) {
		case app_err.ErrUserNotExist:
			return utils.ErrorJson(c, http.StatusNotFound, err.Error())
		case app_err.ErrNoSuchRequest:
			return utils.ErrorJson(c, http.StatusNotFound, err.Error())
    }
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}
	return c.String(http.StatusOK, "OK")
}
