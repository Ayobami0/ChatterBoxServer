package conversation

import (
	"errors"
	"fmt"

	e "github.com/Ayobami0/chatter_box_server/src/errors"
	"github.com/Ayobami0/chatter_box_server/src/model"
	"gorm.io/gorm"
)

func NewGormConversationRepository(db *gorm.DB) *GormConversationRepository {
	return &GormConversationRepository{db}
}

type GormConversationRepository struct {
	DB *gorm.DB
}

func (c *GormConversationRepository) CreateConversation(conversation *model.Conversation) error {
	return c.DB.Create(conversation).Error
}

func (c *GormConversationRepository) DeleteConversation(id string) error {
	err := c.DB.Delete(&model.Conversation{}, "id = ?", id).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return e.ErrNoSuchConversation(id)
		default:
			return err
		}
	}

	return nil
}

func (c *GormConversationRepository) UpdateConversation(update *model.Conversation) error {
	err := c.DB.Model(&model.Conversation{}).Updates(update).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return e.ErrNoSuchConversation(update.ID)
		default:
			return err
		}
	}

	return nil
}

func (c *GormConversationRepository) AddMessage(id string, message model.Message) error {
	var conversation model.Conversation
	err := c.DB.Find(&conversation, "id = ?", id).Error

	if err != nil || conversation.ID == "" {
		switch {
    case conversation.ID == "":
      fallthrough
		case errors.Is(err, gorm.ErrRecordNotFound):
			return e.ErrNoSuchConversation(id)
		default:
			return err
		}
	}

	return c.DB.Model(&conversation).Association("Messages").Append(&message)
}

func (c *GormConversationRepository) DeleteMessage(id string, message model.Message) error {
	var conversation model.Conversation
	err := c.DB.Find(&conversation, "id = ?", id).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return e.ErrNoSuchConversation(id)
		default:
			return err
		}
	}

	return c.DB.Model(&conversation).Association("Messages").Delete(&message)
}

func (c *GormConversationRepository) AcceptRequest(id, rID string) error {
	var conversation model.Conversation
	err := c.DB.Find(&conversation, "id = ?", id).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return e.ErrNoSuchConversation(id)
		default:
			return err
		}
	}
  request := model.Request{
    ID: rID,
    RequestBase: model.RequestBase{
      ConversationId: &id,
    },
  }

  err = c.DB.Model(&conversation).Association("Requests").Delete(&request)

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return e.ErrNoSuchConversation(id)
		default:
			return err
		}
	}

	return c.DB.Model(&conversation).Association("Members").Append(&request.User)
}

func (c *GormConversationRepository) RemoveMember(id, uID string) error {
	var conversation model.Conversation
	err := c.DB.Find(&conversation, "id = ?", id).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return e.ErrNoSuchConversation(id)
		default:
			return err
		}
	}

	user := model.User{
		ID: uID,
	}

	return c.DB.Model(&conversation).Association("Members").Delete(&user)
}

func (c *GormConversationRepository) RejectRequest(id, rID string) error {
	var conversation model.Conversation

	err := c.DB.Find(&conversation, "id = ?", id).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return e.ErrNoSuchConversation(id)
		default:
			return err
		}
	}

  request := model.Request{
    ID: rID,
    RequestBase: model.RequestBase{
      ConversationId: &id,
    },
  }

  err = c.DB.Model(&conversation).Association("Requests").Delete(&request)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return e.ErrNoSuchConversation(id)
		default:
			return err
		}
	}

  return nil
}

func (c *GormConversationRepository) CreateConversationRequest(id string, request model.Request) error {
	var conversation model.Conversation

	err := c.DB.Find(&conversation, "id = ?", id).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return e.ErrNoSuchConversation(id)
		default:
			return err
		}
	}


  return c.DB.Model(&conversation).Association("Requests").Append(&request)
}

func (c *GormConversationRepository) QueryConversations(q string, page, count int, all bool) ([]model.Conversation, error) {
	var conversations []model.Conversation

	res := c.DB.Preload("Members").Limit(count).Offset((page - 1) * count)

	if all {
		res = res.Find(&conversations, "type = ?", "GROUP")
	} else {
		res = res.Find(&conversations, "name LIKE ? AND type = ?", fmt.Sprintf("%%%s%%", q), "GROUP")
	}

	return conversations, res.Error
}

func (c *GormConversationRepository) GetRequests(id string) ([]model.Request, error) {
	var conversation model.Conversation
	var requests []model.Request
	err := c.DB.Find(&conversation, "id = ?", id).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, e.ErrNoSuchConversation(id)
		default:
			return nil, err
		}
	}

	err = c.DB.Model(&conversation).Preload("Conversation").Preload("From").Association("Requests").Find(&requests)
	if err != nil {
		return nil, err
	}

	return requests, nil

}


func (c *GormConversationRepository) GetMessages(id string) ([]model.Message, error) {
	var conversation model.Conversation
  var messages []model.Message
	err := c.DB.Find(&conversation, "id = ?", id).Error

	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, e.ErrNoSuchConversation(id)
		default:
			return nil, err
		}
	}

  err = c.DB.Model(&conversation).Preload("Messages").Association("Messages").Find(&messages)
	if err != nil {
		return nil, err
	}

  return messages, nil
}
