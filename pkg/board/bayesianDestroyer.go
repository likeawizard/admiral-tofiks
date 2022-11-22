package board

import (
	"math/rand"
	"sort"
)

// BayesianDestroyer calculates a super position of all remaining ships in play, then prioritizes the squares with the highest probability of a present ship.
// Takes 57.6 shots to win
type BayesianDestroyer struct {
	board       Board
	aliveShips  []shipConfig
	mode        int
	orientation int
	targets     []int
	ships       []Ship
}

func NewBayesianDestroyer() *BayesianDestroyer {
	r := BayesianDestroyer{}
	r.mode = SEEK
	r.aliveShips = ShipConfig
	r.targets = make([]int, BOARD_SIZE*BOARD_SIZE)
	for i := 0; i < len(r.targets); i++ {
		r.targets[i] = i
	}

	return &r
}

func (r *BayesianDestroyer) FireControl() int {
	heatMap := r.GenerateHeatMap()
	if r.mode == SEEK {
		return r.sqFromHeatMap(*heatMap)
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
			if connectionFn(sq, hullSq) && r.board[sq/BOARD_SIZE][sq%BOARD_SIZE] == 0 {
				found = true
				r.targets[idx] = r.targets[len(r.targets)-1]
				r.targets = r.targets[:len(r.targets)-1]
				return sq
			}
		}
	}
	return sq
}

func (r *BayesianDestroyer) FireStatus(sq, status int) {

	if status == MISS {
		r.board[sq/BOARD_SIZE][sq%BOARD_SIZE] = 1
		return
	}

	r.board[sq/BOARD_SIZE][sq%BOARD_SIZE] = 2
	r.mode = DESTROY
	r.ships = mergeShips(r.ships, sq)
	hitShipIdx := getShipIdx(r.ships, sq)
	r.orientation = shipOrientation(r.ships[hitShipIdx])

	// On sinking remove any adjacent squares as possible targets
	if status == SUNK {
		r.mode = SEEK
		r.orientation = UNKNOWN

		r.targets = eliminateTargetsAroundSunk(r.ships[hitShipIdx], r.targets)
		r.MarkEmptyAroundSunk(r.ships[hitShipIdx])
		r.EliminateAliveTarget(r.ships[hitShipIdx])
		r.ships[hitShipIdx] = r.ships[len(r.ships)-1]
		r.ships = r.ships[:len(r.ships)-1]
	}
}

func (r *BayesianDestroyer) MarkEmptyAroundSunk(ship Ship) {
	for _, sq := range ship.hull {
		adj := Adjacent(sq)
		for _, adjSq := range adj {
			r.board[adjSq/BOARD_SIZE][adjSq%BOARD_SIZE] = 1
		}
	}
}

func (r *BayesianDestroyer) EliminateAliveTarget(ship Ship) {
	for i, aliveTarget := range r.aliveShips {
		if aliveTarget.size == len(ship.hull) {
			r.aliveShips[i].count--
			if r.aliveShips[i].count == 0 {
				r.aliveShips[i] = r.aliveShips[len(r.aliveShips)-1]
				r.aliveShips = r.aliveShips[:len(r.aliveShips)-1]
			}
			return
		}
	}
}

func (r *BayesianDestroyer) GenerateHeatMap() *Board {
	sort.Slice(r.aliveShips, func(i, j int) bool {
		return r.aliveShips[i].size > r.aliveShips[j].size
	})
	heatmap := Board{}
	targetShip := r.aliveShips[0]
	for sq := 0; sq < BOARD_SIZE*BOARD_SIZE; sq++ {
		file, rank := sq/BOARD_SIZE, sq%BOARD_SIZE
		if rank+targetShip.size <= BOARD_SIZE && r.board.isClear(file, rank, targetShip.size, true) {
			for i := 0; i < targetShip.size; i++ {
				heatmap[file][rank+i]++
			}
		}
		if rank+targetShip.size <= BOARD_SIZE && r.board.isClear(file, rank, targetShip.size, false) {
			for i := 0; i < targetShip.size; i++ {
				heatmap[rank+i][file]++
			}
		}
	}
	return &heatmap
}

func (r BayesianDestroyer) sqFromHeatMap(heatMap Board) int {
	type sqScore struct {
		sq     int
		weight int
	}
	squares := make([]sqScore, 0)
	for sq := 0; sq < BOARD_SIZE*BOARD_SIZE; sq++ {
		file, rank := sq/BOARD_SIZE, sq%BOARD_SIZE
		if heatMap[file][rank] != 0 {
			squares = append(squares, sqScore{sq: sq, weight: heatMap[file][rank]})
		}
	}
	sort.Slice(squares, func(i, j int) bool {
		return squares[i].weight > squares[j].weight
	})

	return squares[0].sq
}
