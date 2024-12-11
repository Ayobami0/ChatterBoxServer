package model

type WSMessage struct {
  // type of message. allowed messages - message, update
	Type string      `json:"type"`
  // data
	Data map[string]interface{} `json:"data"`
}
