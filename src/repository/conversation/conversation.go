package conversation

import "github.com/Ayobami0/chatter_box_server/src/model"

type ConversationRepository interface {
	CreateConversation(conv *model.Conversation) error
	DeleteConversation(id string) error
	QueryConversations(q string, page, count int, all bool) ([]model.Conversation, error)
	UpdateConversation(update *model.Conversation) error
}
