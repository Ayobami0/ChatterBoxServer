package handler

import (
	"log"
	"net/http"

	"github.com/Ayobami0/chatter_box_server/src/constant"
	"github.com/Ayobami0/chatter_box_server/src/model"
	"github.com/Ayobami0/chatter_box_server/src/service"
	"github.com/Ayobami0/chatter_box_server/src/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{}

type ChatHandler struct {
	w  *service.WebsocketService
	cS *service.ConversationService
}

func NewChatHandler(w *service.WebsocketService, c *service.ConversationService) ChatHandler {
	return ChatHandler{w, c}
}

func (ch *ChatHandler) ChatConnect(c echo.Context) error {
	u := c.Get("user").(*jwt.Token)
	claim := u.Claims.(*service.UserAccountClaim)

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)

	if err != nil {
		log.Println(claim.UserId, err)
		return utils.ErrorJson(c, http.StatusBadRequest, err.Error())
	}

	ch.w.AddConnection(claim.UserId, conn)
	defer ch.w.RemoveConnection(claim.UserId)

	for {
		var msg model.WSMessage
		err := conn.ReadJSON(&msg)

		if err != nil {
			ch.w.RemoveConnectionWithError(claim.UserId, "message format not supported", websocket.ClosePolicyViolation)
			continue
		}

		switch msg.Type {
		case constant.WS_TYPE_MESSAGE:
			log.Println(msg.Data)
			data := msg.Data
			sender, ok := data["senderId"].(string)
			if !ok {
				ch.w.RemoveConnectionWithError("", "invalid data format for type 'message'", websocket.CloseUnsupportedData)
				if sender == claim.UserId {
					return nil
				}
        continue
			}

			s, ok := data["receiverIds"].([]interface{})
			if !ok {
				ch.w.RemoveConnectionWithError(sender, "invalid data format for type 'message'", websocket.CloseUnsupportedData)
				if sender == claim.UserId {
					return nil
				}
        continue
			}

			cId, ok := data["conversationId"].(string)
			if !ok {
				ch.w.RemoveConnectionWithError(sender, "invalid data format for type 'message'", websocket.CloseUnsupportedData)
				if sender == claim.UserId {
					return nil
				}
        continue
			}

			content, ok := data["content"].(string)
			if !ok {
				ch.w.RemoveConnectionWithError(sender, "invalid data format for type 'message'", websocket.CloseUnsupportedData)
				if sender == claim.UserId {
					return nil
				}
        continue
			}

			attachmentUrl, ok := data["attachmentUrl"].(*string)

			subscribers := make([]string, 0)
			for i := 0; i < len(s) && len(s) < constant.MAX_ALLOWED_GROUP_MEMBERS; i++ { // Allows only a max of 100 subscribers
				if val, ok := s[i].(string); ok {
					if sender != val {
						subscribers = append(subscribers, val)
					}
				}
			}

			subscribers = append(subscribers, sender)

			createdMessage := model.MessageCreate{
				MessageBase: model.MessageBase{
					Content:        content,
					AttachmentUrl:  attachmentUrl,
					ConversationId: cId,
					SenderId:       sender,
				},
			}

			message, err := ch.cS.AddMessage(cId, createdMessage)

			if err != nil {
				ch.w.RemoveConnectionWithError(sender, err.Error(), websocket.CloseNormalClosure)
				if sender == claim.UserId {
					return nil
				}
        continue
			}

			ch.w.BroadcastMessage(message, subscribers...)
		}
	}
}
