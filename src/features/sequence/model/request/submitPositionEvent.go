package request

type ReqWSSubmitPosition struct {
	Round     int     `json:"round"`
	ImageID   int     `json:"imageId"`
	XPosition float64 `json:"xPosition"`
	YPosition float64 `json:"yPosition"`
}
