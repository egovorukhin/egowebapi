package models

import "time"

type Response struct {
	Id       int       `json:"id"`
	Message  string    `json:"message"`
	Datetime time.Time `json:"datetime"`
}

const ModelResponse = "Response"
