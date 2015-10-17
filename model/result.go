package model

type Result int

const (
	Continue Result = iota
	Win
	Loose
)

func (r Result) String() (repr string) {
	switch r {
	case Win:
		repr = "You win!"
	case Loose:
		repr = "You loose :-("
	}
	return
}
