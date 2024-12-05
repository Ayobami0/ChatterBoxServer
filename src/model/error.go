package model

import "fmt"

type Error struct {
  Status string `json:"code"`
  Message string `json:"message"`
}

func NewHttpError(status int, message string) Error {
  return Error{
    Status: fmt.Sprintf("%d", status),
    Message: message,
  }
}
