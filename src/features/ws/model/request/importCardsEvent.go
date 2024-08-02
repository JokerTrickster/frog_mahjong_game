package request

type ReqWSImportCards struct {
	Cards []Card `json:"cards"`
}

type ImportCards struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	State string `json:"state"`
}
