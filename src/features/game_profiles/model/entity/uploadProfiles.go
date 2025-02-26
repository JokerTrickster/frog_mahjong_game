package entity

import "mime/multipart"

type ImageUploadProfileEntity struct {
	Image       *multipart.FileHeader `json:"image"`
	Name        string                `json:"name"`
	TotalCount  int                   `json:"totalCount"`
	Description string                `json:"description"`
}
