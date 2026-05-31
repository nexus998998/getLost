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
	gameSet = "./baseGameSet"
)

// keyboard.Listen
// var binds map[]

var input game.Input

func main() {
	chars := game.Charecters{
		Char1: game.DefaultCharecter,
		Char2: game.DefaultCharecter,
	}
	g, err := game.NewGame(gameSet, chars)
	if err != nil {
		log.Fatal(err)
	}
	// IdleAnimInstance := g.Charecters.Char1.AnimationSet[game.Walking]
	// g.State.P1State.CurrentAnimation = &IdleAnimInstance

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

		g.MakeFrame(g.State).RenderFrame()
		input = game.None // thhis is confusing please change it
		time.Sleep(time.Millisecond * 50)

	}
}
