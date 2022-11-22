package board

const (
	BOARD_SIZE int = 10
)

// Shot status
const (
	MISS int = iota
	HIT
	SUNK
)

// Fire control
const (
	SEEK int = iota
	DESTROY
)

const (
	UNKNOWN int = iota
	HORIZONATL
	VERTICAL
)

type Board [BOARD_SIZE][BOARD_SIZE]int

type Game struct {
	Board Board
	Ships []Ship
	Life  int
	Agent Agent
}

type Ship struct {
	size int
	hull []int
}

type shipConfig struct {
	count int
	size  int
}
