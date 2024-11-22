package entity

import "mime/multipart"

type SaveCardImageGameEntity struct {
	Image *multipart.FileHeader `json:"image"`
}
