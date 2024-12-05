package model

import "time"

type RequestBase struct {
	FromID         string        `json:"fromID"`
  From           User          `json:"from" gorm:"foreignKey:FromID"`
	UserID         *string       `json:"userId,omitempty"`
	User           *User         `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt      time.Time     `json:"createdAt"`
	ConversationId *string       `json:"conversationId,omitempty"`
	Conversation   *Conversation `json:"conversation,omitempty" gorm:"foreignKey:ConversationId"`
}

type RequestCreate struct {
	RequestBase
}

type Request struct {
	RequestBase
	ID string `json:"id" gorm:"primarykey"`
}
