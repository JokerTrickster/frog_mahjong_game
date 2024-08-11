package ws

func CalcPlayTurn(playTurn, playerCount int) int {
	return (playTurn % playerCount) + 1
}
