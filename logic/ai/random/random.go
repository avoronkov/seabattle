package random

import (
	"math/rand"

	"github.com/avoronkov/seabattle/logic/ai"
	"github.com/avoronkov/seabattle/model"
)

func init() {
	maker := func() ai.Interface {
		return RandomAI{}
	}
	ai.AIMakers[""], ai.AIMakers["random"] = maker, maker
}

type RandomAI struct{}

func (e RandomAI) Shoot(b *model.Board) (y, x int) {
	unshoot := 0
	for _, line := range b.Board {
		for _, cell := range line {
			if !cell.Shooted() {
				unshoot++
			}
		}
	}
	n := rand.Intn(unshoot)
	for j, line := range b.Board {
		for i, cell := range line {
			if !cell.Shooted() {
				if n == 0 {
					return j, i
				}
				n--
			}
		}
	}
	panic("Internal error")
}
