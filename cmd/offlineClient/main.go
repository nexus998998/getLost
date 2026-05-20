package main

import (
	"getLost/internals/game"
	"log"
	"os"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
)

var (
	configFilePath = "./config.json"
)

// keyboard.Listen
// var binds map[]

var input game.Input

func main() {
	chars := game.Charecters{
		Char1: game.DefaultCharecter,
		Char2: game.DefaultCharecter,
	}
	g, err := game.NewGame(configFilePath, chars)
	if err != nil {
		log.Fatal(err)
	}
	anim := chars.Char1.AnimationSet[game.Walking]
	g.State.P1State.CurrentAnimation = &anim
	// game loop
	// what's the plan , init the game , make a rendering loop and yeah that's it

	go func() {
		for {
			keyboard.Listen(func(k keys.Key) (bool, error) {
				if k.Code == keys.Left {
					input = game.WalkLeft
					return true, nil
				}
				if k.Code == keys.CtrlC {
					os.Exit(0)
				}

				return true, nil

			})
		}
	}()

	for {
		g.State = g.GenerateNextState(g.State, input)

		f := g.MakeFrame(g.State)
		f.RenderFrame()
		input = game.None // thhis is confusing please change it
		time.Sleep(time.Millisecond * 100)

	}
}
