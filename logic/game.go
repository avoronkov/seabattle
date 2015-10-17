package logic

import (
	"log"

	"github.com/avoronkov/seabattle/logic/ai"
	"github.com/avoronkov/seabattle/model"
	"github.com/avoronkov/seabattle/view"
)

type Game struct {
	my, other *model.Board
	iface     view.Interface
	cpu       ai.Interface
}

func NewGame(viewType, aiType string) (g *Game, err error) {
	g = new(Game)

	g.my = model.NewBoard(config.Height, config.Width, false)
	g.other = model.NewBoard(config.Height, config.Width, true)

	if err = g.my.GenerateShips(); err != nil {
		return
	}
	log.Printf("my field:\n%v\n", *g.my)

	if err = g.other.GenerateShips(); err != nil {
		return
	}
	log.Printf("enemy field:\n%v\n", *g.other)

	if g.iface, err = view.GetView(viewType); err != nil {
		return
	}

	if err = g.iface.SetBoards(g.my, g.other); err != nil {
		return
	}

	if g.cpu, err = ai.Get(aiType); err != nil {
		return
	}

	return
}

func (g *Game) Play() (err error) {
	var (
		y, x int
	)
L:
	for {
		g.iface.Draw()
		if !g.my.HasShipsAlive() {
			g.iface.ShowResult(model.Loose)
			break
		}
		for {
			if y, x, err = g.iface.GetShoot(); err != nil {
				break L
			}
			e := g.other.Shoot(y, x)
			if e == nil {
				break
			}
			if e != model.ErrDoubleShot {
				err = e
				break L
			}
			g.iface.Draw()
		}

		g.iface.Draw()
		if !g.other.HasShipsAlive() {
			g.iface.ShowResult(model.Win)
			break
		}
		y, x = g.cpu.Shoot(g.my)
		if err = g.my.Shoot(y, x); err != nil {
			break
		}
	}
	return
}

func (g *Game) Close() error {
	return g.iface.Close()
}
