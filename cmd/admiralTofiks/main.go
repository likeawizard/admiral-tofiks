package main

import (
	"fmt"

	"github.com/likeawizard/admiral-tofiks/pkg/board"
)

func main() {
	fmt.Println("Admiral Tofiks reporting for duty!")
	shots := 0
	numGames := 10000
	for i := 0; i < numGames; i++ {
		g := board.NewGame()
		g.Agent = board.NewBayesianDestroyer()
		for g.Life > 0 {
			sq := g.Agent.FireControl()
			status := g.Fire(sq)
			g.Agent.FireStatus(sq, status)
			shots++
		}
	}

	fmt.Println("average shots taken", float64(shots)/float64(numGames))
}
