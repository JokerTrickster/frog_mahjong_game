package request

type ReqUpdateGameUsers struct {
	Name      string `json:"name" validate:"omitempty,min=2,max=6" `
	Password  string `json:"password" validate:"omitempty,min=6"`
	ProfileID int    `json:"profileID" validate:"omitempty"`
}
