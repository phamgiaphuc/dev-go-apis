package models

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}
