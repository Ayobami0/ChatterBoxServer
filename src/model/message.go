package model

import "time"

type MessageBase struct {
	Content        string  `json:"content"`
	AttachmentUrl  *string `json:"attachmentUrl"`
	SenderId       string  `json:"senderId"`
	ConversationId string  `json:"conversationId"`
}

type MessageCreate struct {
	MessageBase
}

type Message struct {
	MessageBase
	ID        string    `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"sentAt"`
}
