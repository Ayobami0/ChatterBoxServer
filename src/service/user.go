package service

import (
	"github.com/Ayobami0/chatter_box_server/src/errors"
	"github.com/Ayobami0/chatter_box_server/src/model"
	"github.com/Ayobami0/chatter_box_server/src/repository/user"
	"github.com/google/uuid"
)

func NewUserService(r user.UserRepository) UserService {
	return UserService{r}
}

type UserService struct {
	r user.UserRepository
}

func (u UserService) QueryUsersByUsername(q string, page, count int) ([]model.User, error) {
	if page < 1 {
		page = 1
	}
	if count < 0 {
		count = 10
	}

	users, err := u.r.QueryUsers(q, page, count)

	if err != nil {
		return users, errors.ErrUnexpected("an unexpected error occured")
	}

	return users, nil
}

func (u UserService) GetUserConversations(id string) ([]model.Conversation, error) {
	conversations, err := u.r.GetUserConversations(id)
	if err != nil {
		return nil, err
	}

	return conversations, nil
}

func (u UserService) GetInviteRequests(id string) ([]model.Request, error) {
	requests, err := u.r.GetUserRequests(id)

	return requests, err
}

func (u UserService) SendInviteToUser(fromId, toId string) error {
	rId, err := uuid.NewUUID()

	if err != nil {
		return err
	}

	nReq := model.Request{
		ID: rId.String(),
		RequestBase: model.RequestBase{
			UserID: &toId,
			FromID: fromId,
		},
	}

	return u.r.CreateUserInvite(toId, &nReq)
}

func (u UserService) DeclineRequest(id, reqId string) error {
	return u.r.UpdateUserRequest(id, reqId, nil)
}

func (u UserService) AcceptRequest(id, reqId string) error {
	cId, err := uuid.NewUUID()

	if err != nil {
		return err
	}

	newConv := &model.Conversation{
		ConversationBase: model.ConversationBase{
			Type:      "PERSONAL",
			IsPrivate: true,
		},
		ID: cId.String(),
	}
	return u.r.UpdateUserRequest(id, reqId, newConv)
}
