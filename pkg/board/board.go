package board

import (
	"fmt"
	"math/rand"
	"time"
)

func NewGame() Game {
	g := Game{}
	seed := int64(time.Now().Nanosecond())
	// seed = 677611440
	// fmt.Println("Seed", seed)

	rand.Seed(seed)
	failCount := 0
	// Loop through ship clesses
	for _, ship := range ShipConfig {
		// Loop through ships within class
		for num := 0; num < ship.count; num++ {
			file, rank, horizontal := rand.Intn(BOARD_SIZE), rand.Intn(BOARD_SIZE-ship.size+1), randBool()
			for !g.Board.hasClearance(file, rank, ship.size, horizontal) {
				file, rank, horizontal = rand.Intn(BOARD_SIZE), rand.Intn(BOARD_SIZE-ship.size+1), randBool()
				failCount++
			}
			hull := make([]int, ship.size)
			for i := 0; i < ship.size; i++ {
				g.Life++
				if horizontal {
					g.Board[file][rank+i] = 1
					hull[i] = file*BOARD_SIZE + rank + i
				} else {
					g.Board[rank+i][file] = 1
					hull[i] = (rank+i)*BOARD_SIZE + file
				}
			}
			g.Ships = append(g.Ships, Ship{hull: hull, size: ship.size})
		}
	}

	return g
}

func (b *Board) String() string {
	s := ""
	for file := 0; file < BOARD_SIZE; file++ {
		s += fmt.Sprintf("%2d  ", file+1)
		for rank := 0; rank < BOARD_SIZE; rank++ {
			s += fmt.Sprintf("%d ", b[file][rank])
		}
		s += fmt.Sprintln()
	}
	s += fmt.Sprintln("\n    a b c d e f g h i j")
	return s
}

var ShipConfig = []shipConfig{
	{count: 1, size: 5},
	{count: 1, size: 4},
	{count: 1, size: 3},
	{count: 2, size: 2},
	{count: 2, size: 1},
}

func randBool() bool {
	return rand.Int()%2 == 0
}

func (b *Board) hasClearance(file, rank, size int, horizontal bool) bool {
	clear := func(rank, file int) bool {
		for r := -1; r <= 1; r++ {
			if rank+r >= BOARD_SIZE || rank+r < 0 {
				continue
			}
			for f := -1; f <= 1; f++ {
				if file+f >= BOARD_SIZE || file+f < 0 {
					continue
				}
				if b[rank+r][file+f] == 1 {
					return false
				}
			}
		}
		return true
	}

	for i := 0; i < size; i++ {
		if horizontal {
			if !clear(file, rank+i) {
				return false
			}
		} else {
			if !clear(rank+i, file) {
				return false
			}
		}
	}
	return true
}

func (b *Board) isClear(file, rank, size int, horizontal bool) bool {
	for i := 0; i < size; i++ {
		if horizontal {
			if b[file][rank+i] == 1 {
				return false
			}
		} else {
			if b[rank+i][file] == 1 {
				return false
			}
		}
	}
	return true
}

// Find the ship that has been hit and return true on sinking
func (g *Game) RegisterHit(sq int) bool {
	for shipIdx := 0; shipIdx < len(g.Ships); shipIdx++ {
		for partIdx, part := range g.Ships[shipIdx].hull {
			if part == sq {
				if len(g.Ships[shipIdx].hull) == 1 {
					// fmt.Println("Sunk ship of size", g.Ships[shipIdx].size)
					// fmt.Println("Killing hit", sq)
					g.Ships[shipIdx] = g.Ships[len(g.Ships)-1]
					g.Ships = g.Ships[:len(g.Ships)-1]
					return true
				} else {
					g.Ships[shipIdx].hull[partIdx] = g.Ships[shipIdx].hull[len(g.Ships[shipIdx].hull)-1]
					g.Ships[shipIdx].hull = g.Ships[shipIdx].hull[:len(g.Ships[shipIdx].hull)-1]
					return false
				}
			}
		}
	}
	return false
}

func (g *Game) Fire(sq int) int {
	if g.Board[sq/BOARD_SIZE][sq%BOARD_SIZE] == 1 {
		g.Life--
		if g.RegisterHit(sq) {
			return SUNK
		}
		return HIT
	}
	return MISS
}

func Connected(sq, target int) bool {
	r1, f1 := sq/BOARD_SIZE, sq%BOARD_SIZE
	r2, f2 := target/BOARD_SIZE, target%BOARD_SIZE
	return (f1 == f2 && ((r1 == r2+1) || r1 == r2-1)) || (r1 == r2 && (f1 == f2+1 || f1 == f2-1))
}

func ConnectedV(sq, target int) bool {
	r1, f1 := sq/BOARD_SIZE, sq%BOARD_SIZE
	r2, f2 := target/BOARD_SIZE, target%BOARD_SIZE
	return (f1 == f2 && ((r1 == r2+1) || r1 == r2-1))
}

func ConnectedH(sq, target int) bool {
	r1, f1 := sq/BOARD_SIZE, sq%BOARD_SIZE
	r2, f2 := target/BOARD_SIZE, target%BOARD_SIZE
	return (r1 == r2 && (f1 == f2+1 || f1 == f2-1))
}

func Adjacent(sq int) []int {
	rank, file := sq/BOARD_SIZE, sq%BOARD_SIZE
	adjacent := make([]int, 0)
	for r := -1; r <= 1; r++ {
		if rank+r >= BOARD_SIZE || rank+r < 0 {
			continue
		}
		for f := -1; f <= 1; f++ {
			if file+f >= BOARD_SIZE || file+f < 0 {
				continue
			}
			adjacent = append(adjacent, (rank+r)*BOARD_SIZE+(file+f))
		}
	}
	return adjacent
}
