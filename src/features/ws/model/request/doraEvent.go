package request

type ReqWSDora struct {
	Cards []Card `json:"cards"`
}

type Card struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	State string `json:"state"`
}
