package board

import "sort"

// Agent defines the player. They have two actions available:
// - FireControl provides a target square
// - FireStatus allows the agentg to adjust the FireControl behaviour on MISS/HIT/SUNK responses
type Agent interface {
	FireControl() int
	FireStatus(sq, status int)
}

// Each new disconnected shot will register a new ship. Two disconnected shots can still be hitting the same ship
// When hitting a target see if it can connect existing ships
func mergeShips(ships []Ship, sq int) []Ship {
	newShip := Ship{hull: []int{sq}}
	connected := make([]int, 0)
	for shipIdx, ship := range ships {
		for _, hullElement := range ship.hull {
			if Connected(sq, hullElement) {
				newShip.hull = append(newShip.hull, ship.hull...)
				connected = append(connected, shipIdx)
				break
			}
		}
	}
	ships = append(ships, newShip)
	for _, idx := range connected {
		ships[idx] = ships[len(ships)-1]
		ships = ships[:len(ships)-1]
	}
	return ships
}

func getShipIdx(ships []Ship, sq int) int {
	for shipIdx, ship := range ships {
		for _, hullElement := range ship.hull {
			if hullElement == sq {
				return shipIdx
			}
		}
	}

	return -1
}

func eliminateTargetsAroundSunk(ship Ship, targets []int) []int {
	for _, sq := range ship.hull {
		adj := Adjacent(sq)
		for _, adjSq := range adj {
			for idx, targetSq := range targets {
				if targetSq == adjSq {
					targets[idx] = targets[len(targets)-1]
					targets = targets[:len(targets)-1]
				}
			}
		}
	}
	return targets
}

func shipOrientation(ship Ship) int {
	if len(ship.hull) <= 1 {
		return UNKNOWN
	}
	sort.Slice(ship.hull, func(i, j int) bool {
		return ship.hull[i] < ship.hull[j]
	})
	if ConnectedH(ship.hull[0], ship.hull[1]) {
		return HORIZONATL
	} else {
		return VERTICAL
	}
}
