package curses

import (
	"errors"
	"fmt"
	"log"

	"github.com/avoronkov/seabattle/model"
	"github.com/avoronkov/seabattle/view"

	cur "github.com/rthornton128/goncurses"
)

func init() {
	view.ViewMakers[""] = NewCursesView
	view.ViewMakers["curses"] = NewCursesView
}

type CursesView struct {
	win     *cur.Window
	colored bool

	my, other *model.Board
}

func NewCursesView() (i view.Interface, err error) {
	v := new(CursesView)
	if v.win, err = cur.Init(); err != nil {
		return
	}
	cur.Echo(false)
	cur.Cursor(0)

	if err = v.initCursesColors(); err != nil {
		return
	}

	i = v
	return
}

func (v *CursesView) initCursesColors() (err error) {
	if v.colored = cur.HasColors(); !v.colored {
		return
	}
	if err := cur.StartColor(); err != nil {
		log.Print(err)
	}
	if err := cur.UseDefaultColors(); err != nil {
		log.Print(err)
	}

	colorMap := []int16{
		cur.C_WHITE,
		cur.C_BLUE,
		cur.C_GREEN,
		cur.C_YELLOW,
		cur.C_CYAN,
		cur.C_RED,
	}
	for i, c := range colorMap {

		if e := cur.InitPair(int16(i), c, cur.C_BLACK); e != nil {
			log.Printf("InitPair(%d, %v, %v) failed: %v", i, colorMap[i], cur.C_BLACK, e)
		}
		if e := cur.InitPair(int16(i+len(colorMap)), cur.C_BLACK, colorMap[i]); e != nil {
			log.Printf("InitPair(%d, %v, %v) failed: %v", i, colorMap[i], cur.C_BLACK, e)
		}
	}
	return
}

func (v *CursesView) SetBoards(my, other *model.Board) error {
	v.my, v.other = my, other
	return nil
}

func (v *CursesView) Close() error {
	cur.End()
	return nil
}

func (v *CursesView) printState(s model.State, showShips bool, inverted bool) {
	if v.colored {
		col := int16(s + 1)
		if !showShips && s == model.Ship {
			col = int16(model.Sea + 1)
		}
		if inverted {
			// number of colors in colorMap
			col += 6
		}
		v.win.AttrOn(cur.ColorPair(col))
		defer v.win.AttrOff(cur.ColorPair(col))
	}
	switch s {
	case model.Sea:
		v.win.Print(".")
	case model.Ship:
		if showShips {
			v.win.Print("#")
		} else {
			v.win.Print(".")
		}
	case model.Shot:
		v.win.Print("*")
	case model.Injured, model.Killed:
		v.win.Print("X")
	default:
		panic(fmt.Errorf("Unhandled state: %v", s))
	}
}

func (v *CursesView) Draw() error {
	return v.drawHighliht(-1, -1)
}

func (v *CursesView) drawHighliht(hy, hx int) error {
	header := "   abcdefghij     abcdefghij"
	height, width := v.win.MaxYX()
	startY, startX := height/2-5, (width-len(header))/2
	v.win.Clear()
	v.win.Move(startY, startX)
	v.win.Print(header)
	for j, _ := range v.my.Board {
		y := startY + 2 + j
		v.win.Move(y, startX)
		v.win.Printf("%d", j)
		v.win.Move(y, startX+3)
		for i := 0; i < v.my.Width(); i++ {
			state := v.my.StateOf(j, i)
			v.printState(state, true, false)
		}
		v.win.Move(y, startX+18)
		for i := 0; i < v.other.Width(); i++ {
			state := v.other.StateOf(j, i)
			highlighted := hy == j || hx == i
			v.printState(state, false, highlighted)
		}
		v.win.Printf("  %2d", j)
	}
	v.win.Move(height/2+8, startX)
	v.win.Print(header)
	return nil
}

func (v *CursesView) ShowResult(result model.Result) error {
	str := result.String()
	height, width := v.win.MaxYX()
	startY, startX := height/2+10, (width-len(str))/2
	v.win.Move(startY, startX)
	v.win.Print(str)
	v.win.GetChar()
	return nil
}

func (v *CursesView) GetShoot() (y int, x int, err error) {
	cx, cy := -1, -1
	for {
		switch key := v.win.GetChar(); key {
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j':
			cx = int(key) - int('a')
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			cy = int(key) - int('0')
		case ' ':
			cx, cy = -1, -1
		case 'q':
			err = errors.New("User finished the game")
			return
		}
		if cx >= 0 && cy >= 0 {
			x, y = cx, cy
			return
		}
		v.drawHighliht(cy, cx)
	}
	return
}
