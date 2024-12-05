package service

import (
	"github.com/Ayobami0/chatter_box_server/src/constant"
	"github.com/Ayobami0/chatter_box_server/src/errors"
	"github.com/Ayobami0/chatter_box_server/src/model"
	c "github.com/Ayobami0/chatter_box_server/src/repository/conversation"
	"github.com/google/uuid"
)

func NewConversationService(r c.ConversationRepository) *ConversationService {
	return &ConversationService{r}
}

type ConversationService struct {
	r c.ConversationRepository
}

func (s *ConversationService) CreateConversation(conv model.ConversationCreate, creator *model.User) (*model.Conversation, error) {

	if conv.Name == "" {
		return nil, errors.ErrMissingContent("name")
	}

	if conv.Type == "" {
		return nil, errors.ErrMissingContent("type")
	}

	if _, ok := constant.CONV_TYPE[conv.Type]; !ok {
		return nil, errors.ErrTypeNotSupported(conv.Type)
	}

	id, err := uuid.NewUUID()

	if err != nil {
		return nil, err
	}

	new_conv := model.Conversation{
		ConversationBase: conv.ConversationBase,
		ID:               id.String(),
		Members:          []*model.User{creator},
	}

	err = s.r.CreateConversation(&new_conv)
	if err != nil {
		return nil, err
	}

	return &new_conv, nil
}

func (s *ConversationService) QueryConversationsByName(q string, page, count int) ([]model.Conversation, error) {
	if page < 1 {
		page = 1
	}
	if count < 0 {
		count = 10
	}
	conversations, err := s.r.QueryConversations(q, page, count, q == "")

	if err != nil {
		return conversations, errors.ErrUnexpected("an unexpected error occured")
	}
	return conversations, nil
}
