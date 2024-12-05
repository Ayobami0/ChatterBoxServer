package user

import (
	"errors"
	"fmt"

	_error "github.com/Ayobami0/chatter_box_server/src/errors"
	"github.com/Ayobami0/chatter_box_server/src/model"
	"gorm.io/gorm"
)

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db}
}

type GormUserRepository struct {
	DB *gorm.DB
}

func (u GormUserRepository) GetUser(usernameOrEmail string) *model.User {
	eUser := &model.User{}

	res := u.DB.First(eUser, "username = ? OR email = ?", usernameOrEmail, usernameOrEmail)

	if res.Error != nil {
		return nil
	}

	return eUser
}

func (u GormUserRepository) GetUserConversations(id string) ([]model.Conversation, error) {
	var user model.User
	conversations := []model.Conversation{}

	err := u.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, _error.ErrUserNotExist(id)
		}
		return nil, err
	}

	err = u.DB.Model(&user).Association("Conversations").Find(&conversations)
	if err != nil {
		return nil, err
	}

	return conversations, nil
}

func (u GormUserRepository) CreateUser(user model.User) error {
	return u.DB.Create(user).Error
}

func (u GormUserRepository) QueryUsers(q string, page, count int) ([]model.User, error) {
	users := []model.User{}

	res := u.DB.Select("Username", "ProfileUrl", "ID").Limit(count).Offset((page-1)*count).Find(&users, "username LIKE ?", fmt.Sprintf("%%%s%%", q))

	return users, res.Error
}

func (u GormUserRepository) GetUserById(id string) (*model.User, error) {
	var user model.User

	err := u.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, _error.ErrUserNotExist(id)
		}
		return nil, err
	}

	return &user, nil
}

func (u GormUserRepository) GetUserRequests(id string) ([]model.Request, error) {
	var user model.User
	requests := []model.Request{}

	err := u.DB.Where("id = ?", id).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, _error.ErrUserNotExist(id)
		}
		return nil, err
	}

	err = u.DB.Model(&user).Preload("User").Preload("From").Association("Invitations").Find(&requests)
	if err != nil {
		return nil, err
	}

	return requests, nil
}

func (u GormUserRepository) CreateUserInvite(toId string, request *model.Request) error {
	var user model.User

	err := u.DB.Where("id = ?", toId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return _error.ErrUserNotExist(toId)
		}
		return err
	}

	err = u.DB.Model(&user).Association("Invitations").Append(request)
	fmt.Println(err)

	return err
}

func (u GormUserRepository) UpdateUserRequest(userId, requestId string, conversation *model.Conversation) error {
	var user model.User
	var request model.Request

	res := u.DB.Where("id = ?", userId).First(&user)

	err := res.Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return _error.ErrUserNotExist(userId)
		}
		return err
	}

	err = u.DB.First(&request, "id = ?", requestId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return _error.ErrRequestNotExist
		}
		return err
	}

	err = u.DB.Model(&user).Association("Invitations").Delete(&request)
	if err != nil {
		return err
	}

	if conversation != nil {
		conversation.Members = []*model.User{
			{ID: userId},
			{ID: request.FromID},
		}
		return u.DB.Create(conversation).Error
	}

	return nil
}
