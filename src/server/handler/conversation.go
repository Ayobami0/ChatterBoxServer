package handler

import (
	"net/http"
	"strconv"

	app_err "github.com/Ayobami0/chatter_box_server/src/errors"
	"github.com/Ayobami0/chatter_box_server/src/model"
	"github.com/Ayobami0/chatter_box_server/src/service"
	"github.com/Ayobami0/chatter_box_server/src/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ConversationHandler struct {
	s *service.ConversationService
}

func NewConversationHandler(service *service.ConversationService) ConversationHandler {
	return ConversationHandler{s: service}
}

func (ch *ConversationHandler) ConversationCreate(c echo.Context) error {
  u := c.Get("user").(*jwt.Token)
  claim := u.Claims.(*service.UserAccountClaim)

  user := model.User{
    ID: claim.UserId,
  }
	nConversation := new(model.ConversationCreate)

	if err := c.Bind(nConversation); err != nil {
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}

	conversation, err := ch.s.CreateConversation(*nConversation, &user)

	if err != nil {
		switch err.(type) {
		case app_err.ErrMissingContent:
			return utils.ErrorJson(c, http.StatusBadRequest, err.Error())
		case app_err.ErrTypeNotSupported:
			return utils.ErrorJson(c, http.StatusBadRequest, err.Error())
		}
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}

	return c.JSON(http.StatusCreated, conversation)
}

func (ch *ConversationHandler) ConversationGet(c echo.Context) error {
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

	conversations, err := ch.s.QueryConversationsByName(search, page, count)

	if err != nil {
		return utils.ErrorJson(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, conversations)
}

func (ch *ConversationHandler) ConversationDelete(c echo.Context) error {
  id := c.Param("id")

  err := ch.s.DeleteConversation(id)

	if err != nil {
		switch err.(type) {
		case app_err.ErrNoSuchConversation:
			return utils.ErrorJson(c, http.StatusBadRequest, err.Error())
    }
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}
	return c.String(http.StatusOK, "")
}

func (ch *ConversationHandler) ConversationsRequestsGet(c echo.Context) error {
  id := c.Param("id")

  convs, err := ch.s.ConversationRequests(id)

	if err != nil {
		switch err.(type) {
		case app_err.ErrNoSuchConversation:
			return utils.ErrorJson(c, http.StatusBadRequest, err.Error())
    }
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}
	return c.JSON(http.StatusOK, convs)
}

func (ch *ConversationHandler) ConversationJoin(c echo.Context) error {
  user := c.Get("user").(*jwt.Token)
  claim := user.Claims.(*service.UserAccountClaim)

  id := c.Param("id")

  err := ch.s.ConversationJoin(id, claim.UserId)

	if err != nil {
		switch err.(type) {
		case app_err.ErrNoSuchConversation:
			return utils.ErrorJson(c, http.StatusBadRequest, err.Error())
    }
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}

	return c.String(http.StatusOK, "")
}

func (ch *ConversationHandler) ConversationReject(c echo.Context) error {
  cID := c.Param("cid")
  rID := c.Param("rid")

  err := ch.s.ConversationRequestReject(cID, rID)

	if err != nil {
		switch err.(type) {
		case app_err.ErrNoSuchConversation:
			return utils.ErrorJson(c, http.StatusBadRequest, err.Error())
    }
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}

	return c.String(http.StatusOK, "")
}

func (ch *ConversationHandler) ConversationAccept(c echo.Context) error {
  cID := c.Param("cid")
  rID := c.Param("rid")

  err := ch.s.ConversationRequestAccept(cID, rID)

	if err != nil {
		switch err.(type) {
		case app_err.ErrNoSuchConversation:
			return utils.ErrorJson(c, http.StatusBadRequest, err.Error())
    }
		return utils.ErrorJson(c, http.StatusBadRequest, "bad request")
	}

	return c.String(http.StatusOK, "")
}
