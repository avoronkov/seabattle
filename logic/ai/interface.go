package ai

import (
	"fmt"

	"github.com/avoronkov/seabattle/model"
)

type Interface interface {
	Shoot(b *model.Board) (y, x int)
}

var AIMakers = map[string]func() Interface{}

func Get(aiType string) (in Interface, err error) {
	var (
		ok    bool
		maker func() Interface
	)

	if maker, ok = AIMakers[aiType]; !ok {
		err = fmt.Errorf("Unknown AI type: %q", aiType)
		return
	}
	in = maker()
	return
}
