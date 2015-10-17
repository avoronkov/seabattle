package view

import (
	"fmt"

	"github.com/avoronkov/seabattle/model"
)

type Interface interface {
	SetBoards(my, other *model.Board) error
	Draw() error
	ShowResult(model.Result) error
	GetShoot() (int, int, error)
	Close() error
}

var ViewMakers = map[string]func() (Interface, error){}

func GetView(v string) (in Interface, err error) {
	var (
		ok    bool
		maker func() (Interface, error)
	)
	if maker, ok = ViewMakers[v]; !ok {
		err = fmt.Errorf("Unknown view interface: %q", v)
		return
	}
	in, err = maker()
	return
}
