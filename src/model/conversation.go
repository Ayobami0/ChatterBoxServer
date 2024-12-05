package model

import "time"

type ConversationBase struct {
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	ImageUrl  *string `json:"imageUrl"`
	Desc      *string `json:"desc"`
	IsPrivate bool    `json:"isPrivate"`
}

type ConversationCreate struct {
	ConversationBase
}

type Conversation struct {
	ConversationBase
	ID        string     `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	Messages  []*Message `json:"messages"`
	Members   []*User    `json:"members" gorm:"many2many:conversation_members;"`
}
