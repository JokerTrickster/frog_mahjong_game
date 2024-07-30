package request

type ReqV02Join struct {
	Tkn      string `query:"tkn"`
	Password string `json:"password,omitempty"`
}
