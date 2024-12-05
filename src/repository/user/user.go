package user

import "github.com/Ayobami0/chatter_box_server/src/model"

type UserRepository interface {
	GetUser(usernameOrEmail string) *model.User
	CreateUser(user model.User) error
  QueryUsers(q string, page, count int) ([]model.User, error)
  GetUserConversations(id string) ([]model.Conversation, error)
  GetUserById(id string) (*model.User, error)
  GetUserRequests(id string) ([]model.Request, error)
  CreateUserInvite(toId string, request *model.Request) error
  UpdateUserRequest(userId, requestId string, conversation *model.Conversation) error
}
