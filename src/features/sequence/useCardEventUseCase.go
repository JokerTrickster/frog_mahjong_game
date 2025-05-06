package sequence

import "main/utils/db/mysql"

func CreateNextKingIndex(kingIndex int, cardInfo *mysql.SequenceCards) int {
	// Convert 1D index to 2D coordinates (9x9 grid)
	currentX := kingIndex % 9
	currentY := kingIndex / 9

	// Direction mapping
	// 0: up-left, 1: up, 2: up-right
	// 3: left, 4: right
	// 5: down-left, 6: down, 7: down-right
	direction := cardInfo.Direction

	// Calculate new position based on direction
	newX := currentX
	newY := currentY

	switch direction {
	case 0: // up-left
		newX--
		newY--
	case 1: // up
		newY--
	case 2: // up-right
		newX++
		newY--
	case 3: // left
		newX--
	case 4: // right
		newX++
	case 5: // down-left
		newX--
		newY++
	case 6: // down
		newY++
	case 7: // down-right
		newX++
		newY++
	}

	// Check boundaries (0 to 8)
	if newX < 0 {
		newX = 0
	} else if newX > 8 {
		newX = 8
	}

	if newY < 0 {
		newY = 0
	} else if newY > 8 {
		newY = 8
	}

	// Convert back to 1D index
	nextKingIndex := newY*9 + newX

	return nextKingIndex
}
