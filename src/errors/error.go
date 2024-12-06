package errors

import (
	"errors"
	"fmt"

	"github.com/Ayobami0/chatter_box_server/src/constant"
)

// Returned if a required parameter is missing from a json request
type ErrMissingContent string

func (e ErrMissingContent) Error() string {
	return fmt.Sprintf("a required parameter is missing from the request body (%s)", string(e))
}

// Returned if an unsupported conversation type is send in json request
type ErrTypeNotSupported string

func (e ErrTypeNotSupported) Error() string {
	keys := make([]string, 0, len(constant.CONV_TYPE))
	for k := range constant.CONV_TYPE {
		keys = append(keys, k)
	}

	return fmt.Sprintf("unsupported type (%s). allowed types %v", string(e), keys)
}

// Returned when something ErrUnexpected occurs from the server
type ErrUnexpected string

func (e ErrUnexpected) Error() string {
	return string(e)
}

type ErrUserExist struct {
	Username string
	Email    string
}

func (e ErrUserExist) Error() string {
	return fmt.Sprintf("user with email `%s` or username `%s` already exist", e.Email, e.Username)
}

type ErrUserNotExist string

func (e ErrUserNotExist) Error() string {
	return fmt.Sprintf("user `%s` does not exist", string(e))
}

type ErrPasswordIncorrect struct{}

func (e ErrPasswordIncorrect) Error() string {
	return "user password is incorrect"
}

type ErrNoSuchConversation string

func (e ErrNoSuchConversation) Error() string {
	return string(e)
}

type ErrNoSuchRequest string

func (e ErrNoSuchRequest) Error() string {
	return string(e)
}

var (
  ErrRequestNotExist = errors.New("request does't exist")
  ErrConvertationPrivate = errors.New("cannot join private conversation")
  ErrMemberNotAdmin = errors.New("member is not an admin")
)
