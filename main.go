package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"time"

	"github.com/avoronkov/seabattle/logic"
	_ "github.com/avoronkov/seabattle/logic/ai/random"
	_ "github.com/avoronkov/seabattle/view/curses"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic: %v : %s", r, debug.Stack())
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error occured: %v\n", err)
			os.Exit(1)
		}
	}()
	game, err := logic.NewGame("", "")
	if err != nil {
		return
	}
	defer func() {
		err = game.Close()
	}()
	err = game.Play()
}
