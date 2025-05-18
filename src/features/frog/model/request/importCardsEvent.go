package request

type ReqWSImportCards struct {
	Cards []ImportCards `json:"cards"`
}

type ImportCards struct {
	CardID uint `json:"cardID"`
}
