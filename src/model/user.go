package model

type UserBase struct {
	Username   string  `json:"username,omitempty" gorm:"unique"`
	Email      string  `json:"email,omitempty" gorm:"unique"`
	ProfileUrl *string `json:"profileUrl"`
	Password   string  `json:"password,omitempty"`
}

type UserCreate struct {
	UserBase
}

type UserLogin struct {
	Password        string `json:"password,omitempty"`
	EmailOrUsername string `json:"emailOrUsername,omitempty"`
}

type User struct {
	UserBase
	ID            string         `json:"id,omitempty" gorm:"primarykey"`
  Password      string         `json:"-"`
	Conversations []*Conversation `json:"conversations,omitempty" gorm:"many2many:conversation_members;"`
	Invitations   []*Request      `json:"invitations,omitempty"`
}
