package board

import "math/rand"

// The random shooter simply fires at squares that have not yet been shot at.
// On average takes 95.7 shots to win
type RandomShooter struct {
	targets []int
}

func NewRandomShooter() *RandomShooter {
	r := RandomShooter{}
	r.targets = make([]int, BOARD_SIZE*BOARD_SIZE)
	for i := 0; i < len(r.targets); i++ {
		r.targets[i] = i
	}
	return &r
}

func (r *RandomShooter) FireControl() int {
	idx := rand.Intn(len(r.targets))
	sq := r.targets[idx]
	r.targets[idx] = r.targets[len(r.targets)-1]
	r.targets = r.targets[:len(r.targets)-1]
	return sq
}

func (r *RandomShooter) FireStatus(sq, status int) {
	// random doesn't care
}
