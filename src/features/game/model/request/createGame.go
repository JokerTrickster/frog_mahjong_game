package request

import "mime/multipart"

type ReqCreateMission struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Image       *multipart.FileHeader `json:"image"`
}
