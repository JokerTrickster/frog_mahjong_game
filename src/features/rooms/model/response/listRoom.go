package response

type ResListRoom struct {
	Total int        `json:"total"`
	Rooms []ListRoom `json:"rooms"`
}

type ListRoom struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Created      int64  `json:"created"`
	CurrentCount int    `json:"currentCount"`
	MaxCount     int    `json:"maxCount"`
	MinCount     int    `json:"minCount"`
	State        string `json:"state"`
	Password     string `json:"password,omitempty"`
	OwnerID      int    `json:"ownerID"`
}
