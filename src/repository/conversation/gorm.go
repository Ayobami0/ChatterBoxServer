package conversation

import (
	"fmt"

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
	return c.DB.Select("Messages").Delete(&model.Conversation{}, id).Error
}

func (c *GormConversationRepository) UpdateConversation(update *model.Conversation) error {
	return c.DB.Model(&model.Conversation{}).Updates(update).Error
}

func (c *GormConversationRepository) AddMessage(id string, message model.Message) error {
	return c.DB.Model(&model.Conversation{}).Where(&model.Conversation{ID: id}).Association("Messages").Append(&message)
}

func (c *GormConversationRepository) DeleteMessage(id string, message model.Message) error {
	return c.DB.Model(&model.Conversation{}).Where(&model.Conversation{ID: id}).Association("Messages").Delete(&message)
}

func (c *GormConversationRepository) AddRequest(id string, request model.Request) error {
	return c.DB.Model(&model.Conversation{}).Where(&model.Conversation{ID: id}).Association("Requests").Append(&request)
}

func (c *GormConversationRepository) AcceptRequest(id string, request model.Request) error {
	model := c.DB.Model(&model.Conversation{}).Where(&model.Conversation{ID: id})
	if err := model.Association("Messages").Delete(&request); err != nil {
		return err
	}

	return model.Association("Members").Append(&request.User)
}

func (c *GormConversationRepository) RejectRequest(id string, request model.Request) error {
	return c.DB.Model(&model.Conversation{}).Association("Messages").Delete(&request)
}

func (c *GormConversationRepository) QueryConversations(q string, page, count int, all bool) ([]model.Conversation, error) {
	var conversations []model.Conversation

	res := c.DB.Preload("Members").Preload("Messages").Limit(count).Offset((page - 1) * count)

	if all {
		res = res.Find(&conversations)
	} else {
		res = res.Find(&conversations, "name LIKE ?", fmt.Sprintf("%%%s%%", q))
	}

	return conversations, res.Error
}
