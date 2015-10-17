package model

import (
	"fmt"
	"math/rand"
)

type State int

const (
	Sea State = iota
	Ship
	Shot
	Injured
	Killed
)

func (s State) String() string {
	switch s {
	case Sea:
		return "."
	case Ship:
		return "#"
	case Shot:
		return "*"
	case Injured:
		return "X"
	default:
		panic(fmt.Errorf("Unknown state: %d", s))
	}
}

func (s State) Shooted() bool {
	return s >= Shot
}

type Board struct {
	Board         [][]State
	height, width int
	// whether ships are hidden
	hidden bool
	ships  map[point]*ship
}

func NewBoard(h, w int, hidden bool) *Board {
	b := new(Board)
	b.Board = make([][]State, h)
	for i := 0; i < h; i++ {
		b.Board[i] = make([]State, w)
	}
	b.height, b.width = h, w
	b.hidden = hidden
	b.ships = make(map[point]*ship)
	return b
}

func (b Board) String() (repr string) {
	for _, line := range b.Board {
		for _, c := range line {
			repr += c.String()
		}
		repr += "\n"
	}
	return
}

func (b *Board) GenerateShips() error {
	for l := 4; l > 0; l-- {
		for i := 0; i < 5-l; i++ {
			b.generateShip(l)
		}
	}
	return nil
}

func (b *Board) generateShip(shipLen int) {
L:
	for {
		if rand.Intn(2) > 0 {
			y, x := rand.Intn(b.height), rand.Intn(b.width-shipLen+1)
			for j := -1; j <= 1; j++ {
				if y+j < 0 || y+j >= b.height {
					continue
				}
				for i := -1; i <= shipLen; i++ {
					if x+i < 0 || x+i >= b.width {
						continue
					}
					if b.Board[y+j][x+i] != Sea {
						continue L
					}
				}
			}
			sh := new(ship)
			for i := 0; i < shipLen; i++ {
				b.Board[y][x+i] = Ship
				sh.points = append(sh.points, &b.Board[y][x+i])
				b.ships[point{y: y, x: x + i}] = sh
			}
		} else {
			y, x := rand.Intn(b.height-shipLen+1), rand.Intn(b.width)
			for j := -1; j <= 1; j++ {
				if x+j < 0 || x+j >= b.width {
					continue
				}
				for i := -1; i <= shipLen; i++ {
					if y+i < 0 || y+i >= b.height {
						continue
					}
					if b.Board[y+i][x+j] != Sea {
						continue L
					}
				}
			}
			sh := new(ship)
			for i := 0; i < shipLen; i++ {
				b.Board[y+i][x] = Ship
				sh.points = append(sh.points, &b.Board[y+i][x])
				b.ships[point{y: y + i, x: x}] = sh
			}
		}
		return
	}
}

func (b *Board) Shoot(y, x int) error {
	var ns State
	switch state := b.Board[y][x]; state {
	case Sea:
		ns = Shot
	case Ship:
		ns = Injured
	default:
		return fmt.Errorf("Already shoot there: (%d, %d)", x, y)
	}
	b.Board[y][x] = ns
	return nil
}

func (b Board) HasShipsAlive() bool {
	for _, line := range b.Board {
		for _, cell := range line {
			if cell == Ship {
				return true
			}
		}
	}
	return false
}

func (b Board) Height() int {
	return b.height
}

func (b Board) Width() int {
	return b.width
}

func (b Board) StateOf(y, x int) (st State) {
	st = b.Board[y][x]
	if st == Ship && b.hidden {
		st = Sea
	}
	if st >= Injured {
		sh := *b.ships[point{y: y, x: x}]
		if sh.injured() {
			st = Injured
		}
		if sh.killed() {
			st = Killed
		}
	}
	return
}
