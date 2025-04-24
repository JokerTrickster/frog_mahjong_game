package response

type ResSlimeWarGetsCardBoardGame struct {
	CardList []Card `json:"cardList"`
}

type Card struct {
	ID        int    `json:"id"`
	Direction int    `json:"direction"`
	ImageUrl  string `json:"imageUrl"`
	Move      int    `json:"move"`
}
