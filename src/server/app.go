package server

import (
	"fmt"

	"github.com/Ayobami0/chatter_box_server/src/constant"
	"github.com/Ayobami0/chatter_box_server/src/repository/conversation"
	"github.com/Ayobami0/chatter_box_server/src/repository/user"
	"github.com/Ayobami0/chatter_box_server/src/server/handler"
	"github.com/Ayobami0/chatter_box_server/src/server/middleware"
	"github.com/Ayobami0/chatter_box_server/src/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func App(config AppConfig) (*echo.Echo, error) {
	// Initialization
	_d, err := config.DB("gorm")
	if err != nil {
		return nil, err
	}
	gormDB, ok := _d.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("invalid database assigned. expected '%T' got '%T'", gorm.DB{}, _d)
	}

	// Repositories
	c_r := conversation.NewGormConversationRepository(gormDB)
	u_r := user.NewGormUserRepository(gormDB)

	// Services
	c_s := service.NewConversationService(c_r)
	a_s := service.NewAuthService(u_r, "secret")
	u_s := service.NewUserService(u_r)
  w_s := service.NewWebsocketService()

	// Handlers
	c_handler := handler.NewConversationHandler(c_s)
	a_handler := handler.NewAuthHandler(a_s)
	u_handler := handler.NewUserHandler(u_s)
  w_handler := handler.NewChatHandler(w_s, c_s)

	e := echo.New()
	auth := e.Group("")
	chat := e.Group("")

	// Middleware

	chat.Use(middleware.WSAuth, a_s.JWTConfig())
	auth.Use(a_s.JWTConfig())

	// Routes
	e.POST(constant.LOGIN_ENDPOINT, a_handler.UserLogin)
	e.PATCH(constant.LOGOUT_ENDPOINT, a_handler.UserLogout)
	e.POST(constant.SIGNUP_ENDPOINT, a_handler.UserCreate)

  chat.GET(constant.WEBSOCKET_ENDPOINT, w_handler.ChatConnect)

	auth.POST(constant.CONVERSATION_CREATE_ENDPOINT, c_handler.ConversationCreate)
	auth.GET(constant.CONVERSATION_REQUEST_ENDPOINT, c_handler.ConversationsRequestsGet)
	auth.PUT(constant.CONVERSATION_JOIN_ENDPOINT, c_handler.ConversationJoin)
	auth.DELETE(constant.CONVERSATION_DELETE_ENDPOINT, c_handler.ConversationDelete)
	auth.PATCH(constant.CONVERSATION_REQUEST_ACCEPT, c_handler.ConversationAccept)
	auth.PATCH(constant.CONVERSATION_REQUEST_REJECT, c_handler.ConversationReject)
	auth.GET(constant.CONVERSATION_GET_ENDPOINT, c_handler.ConversationGet)

	auth.GET(constant.USER_GET_ENDPOINT, u_handler.UsersGet)
	auth.GET(constant.USER_CONVERSATONS_ENDPOINT, u_handler.UserConversationsGet)
	auth.GET(constant.USER_PROFILE_ENDPOINT, u_handler.UsersMeGet)
	auth.GET(constant.USER_REQUEST_ENDPOINT, u_handler.UserRequestsGet)
	auth.PATCH(constant.USER_REQUEST_ACCEPT_ENDPOINT, u_handler.UserRequestsAccept)
	auth.PATCH(constant.USER_REQUEST_REJECT_ENDPOINT, u_handler.UserRequestsReject)
	auth.PUT(constant.USER_REQUEST_INVITE_ENDPOINT, u_handler.UserInvite)


	return e, nil
}
