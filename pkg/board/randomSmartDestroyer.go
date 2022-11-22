package board

import (
	"math/rand"
)

// RandomSmartDestroyer tries to figure out the orientation of the ship it is currently destroying
// Takes 57.6 shots to win
type RandomSmartDestroyer struct {
	mode        int
	orientation int
	targets     []int
	ships       []Ship
}

func NewRandomSmartDestroyer() *RandomSmartDestroyer {
	r := RandomSmartDestroyer{}
	r.mode = SEEK
	r.targets = make([]int, BOARD_SIZE*BOARD_SIZE)
	for i := 0; i < len(r.targets); i++ {
		r.targets[i] = i
	}
	return &r
}

func (r *RandomSmartDestroyer) FireControl() int {
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
			connectionFn := Connected
			switch r.orientation {
			case VERTICAL:
				connectionFn = ConnectedV
			case HORIZONATL:
				connectionFn = ConnectedH
			}
			if connectionFn(sq, hullSq) {
				found = true
				r.targets[idx] = r.targets[len(r.targets)-1]
				r.targets = r.targets[:len(r.targets)-1]
				return sq
			}
		}
	}
	return sq
}

func (r *RandomSmartDestroyer) FireStatus(sq, status int) {
	if status == MISS {
		return
	}
	r.mode = DESTROY
	r.ships = mergeShips(r.ships, sq)
	hitShipIdx := getShipIdx(r.ships, sq)
	r.orientation = shipOrientation(r.ships[hitShipIdx])

	// On sinking remove any adjacent squares as possible targets
	if status == SUNK {
		r.mode = SEEK
		r.orientation = UNKNOWN
		r.targets = eliminateTargetsAroundSunk(r.ships[hitShipIdx], r.targets)
		r.ships[hitShipIdx] = r.ships[len(r.ships)-1]
		r.ships = r.ships[:len(r.ships)-1]
	}
}
