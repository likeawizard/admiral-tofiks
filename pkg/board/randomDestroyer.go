package board

import (
	"math/rand"
)

// RandomDestroyer behaves much like the RandomEliminator but on scoring a hit it will target connected squares to sink it's target
// Takes 64.3 shots to win
type RandomDestroyer struct {
	mode    int
	targets []int
	ships   []Ship
}

func NewRandomDestroyer() *RandomDestroyer {
	r := RandomDestroyer{}
	r.mode = SEEK
	r.targets = make([]int, BOARD_SIZE*BOARD_SIZE)
	for i := 0; i < len(r.targets); i++ {
		r.targets[i] = i
	}
	return &r
}

func (r *RandomDestroyer) FireControl() int {
	if r.mode == SEEK {
		idx := rand.Intn(len(r.targets))
		sq := r.targets[idx]
		r.targets[idx] = r.targets[len(r.targets)-1]
		r.targets = r.targets[:len(r.targets)-1]
		return sq
	}

	found := false
	var sq int
	for !found {
		idx := rand.Intn(len(r.targets))
		sq = r.targets[idx]
		for _, hullSq := range r.ships[0].hull {
			if Connected(sq, hullSq) {
				found = true
				r.targets[idx] = r.targets[len(r.targets)-1]
				r.targets = r.targets[:len(r.targets)-1]
				return sq
			}
		}
	}
	return sq
}

func (r *RandomDestroyer) FireStatus(sq, status int) {
	if status == MISS {
		return
	}
	r.mode = DESTROY
	r.ships = mergeShips(r.ships, sq)
	hitShipIdx := getShipIdx(r.ships, sq)

	// On sinking remove any adjacent squares as possible targets
	if status == SUNK {
		r.mode = SEEK
		r.targets = eliminateTargetsAroundSunk(r.ships[hitShipIdx], r.targets)
		r.ships[hitShipIdx] = r.ships[len(r.ships)-1]
		r.ships = r.ships[:len(r.ships)-1]
	}
}
