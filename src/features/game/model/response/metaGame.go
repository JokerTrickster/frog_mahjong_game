package response

type ResMetaGame struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	ID     uint   `json:"id"`
	Reason string `json:"reason"`
}
