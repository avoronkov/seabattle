package model

type point struct {
	y, x int
}

type ship struct {
	points []*State
}

func (s ship) injured() bool {
	var alives, deads bool
	for _, st := range s.points {
		if *st == Ship {
			alives = true
		}
		if *st == Injured || *st == Killed {
			deads = true
		}
	}
	return alives && deads
}

func (s ship) killed() bool {
	for _, st := range s.points {
		if *st == Ship {
			return false
		}
	}
	return true
}
