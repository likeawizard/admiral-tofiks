package board

import (
	"math/rand"
)

// Random eliminator still shoots randomly. However, it recognizes when a ship has been sunk and marks adjacent squares as ineligible for other ships
// When scoring a hit it registers a new ship or in case it has hit a connected square it adds it to an existing target.
// Takes 81.8 shots to win
type RandomEliminator struct {
	targets []int
	ships   []Ship
}

func NewRandomEliminator() *RandomEliminator {
	r := RandomEliminator{}
	r.targets = make([]int, BOARD_SIZE*BOARD_SIZE)
	for i := 0; i < len(r.targets); i++ {
		r.targets[i] = i
	}
	return &r
}

func (r *RandomEliminator) FireControl() int {
	idx := rand.Intn(len(r.targets))
	sq := r.targets[idx]
	r.targets[idx] = r.targets[len(r.targets)-1]
	r.targets = r.targets[:len(r.targets)-1]
	return sq
}

func (r *RandomEliminator) FireStatus(sq, status int) {
	if status == MISS {
		return
	}
	r.ships = mergeShips(r.ships, sq)
	hitShipIdx := getShipIdx(r.ships, sq)

	// On sinking remove any adjacent squares as possible targets
	if status == SUNK {
		r.targets = eliminateTargetsAroundSunk(r.ships[hitShipIdx], r.targets)
		r.ships[hitShipIdx] = r.ships[len(r.ships)-1]
		r.ships = r.ships[:len(r.ships)-1]
	}
}
