package conversation

import "github.com/Ayobami0/chatter_box_server/src/model"

type ConversationRepository interface {
	CreateConversation(conv *model.Conversation) error
	DeleteConversation(id string) error
	QueryConversations(q string, page, count int, all bool) ([]model.Conversation, error)
	UpdateConversation(update *model.Conversation) error
  RejectRequest(id, rID string) error
  AcceptRequest(id, rID string) error
  RemoveMember(id, uID string) error
  DeleteMessage(id string, message model.Message) error
  AddMessage(id string, message model.Message) error
  GetRequests(id string) ([]model.Request, error)
  GetMessages(id string) ([]model.Message, error)
  CreateConversationRequest(id string, request model.Request) error
}
