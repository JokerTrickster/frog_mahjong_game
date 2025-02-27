package entity

import "mime/multipart"

type FindItImageGameEntity struct {
	Image *multipart.FileHeader `json:"image"`
}
