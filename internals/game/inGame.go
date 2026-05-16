package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

/*
humble to do list :
1- please order structure of the code to make it more readable

*/

// how do i tell the renderer what's the current frame ?
// simple , currentSprite

type Action int

const (
	Idle Action = iota
	Walking
	Ability1
	Ability2
	// etc ....
)

type Direction int

const (
	Right Direction = iota
	Left
)

// listening for input
// when an input is inputted it will trigger a changed state or if there was input blocking (like being in a stun) it will do nothing
// it will change the state as a whole
// for now our mission is to make the state thing and make an MP state constructer

var (
	// use this for colors and special asciis
	ClearAscii = "\033[2J"
	green      = "\033[32m"
	yellow     = "\033[33m"
	blue       = "\033[34m"
	red        = "\033[31m"
	reset      = "\033[0m"
	white      = "\033[37m"
)

var (
	// if players state would be to saperate this would become two variables , ugly
	DefaultState = State{
		P1State: PlayerState{
			Position: Point{10, 10},
		},
		P2State: PlayerState{
			Position: Point{30, 10},
		},
	}

	DefaultCharecter = Charecter{
		Health: 100,
		Speed:  2,
		Hitbox: Hitbox{Point{-1, -1}, Point{1, 1}},
		// optimize animations coding , maybe make a constructer and extract it from an outer file
		AnimationSet: AnimationSet{
			Idle: Animation{
				Panels: [][]string{
					{
						"*-*",
						"-|-",
						" | ",
					},
				},
				RepeatingAnimation: true,
				CurrentFrame:       0,
				AnimationFrames:    2,
			},

			Walking: Animation{
				Panels: [][]string{
					{
						"*-*",
						"-|-",
						" | ",
					},
					{
						"*-*",
						"-|-",
						"/ \\",
					},
				},
				RepeatingAnimation: true,
				CurrentFrame:       0,
				AnimationFrames:    2,
			},
		},
	}

	AnimationIsOver = errors.New("animation is over")
)

// how do i tell the renderer what the charecter is currently doing ?
// ok what sets the actual frame then ?
// how to go for state handlilng ? something like action in state where you have different options for what it could be ?

type Hitbox struct {
	UpperLeft   Point
	BottomRight Point
}

// this is a problem for later , just a sketch
// first layer is the panel layer , the other layer is for 2d rendering
type Animation struct {
	Panels             [][]string
	CurrentFrame       int
	AnimationFrames    int
	RepeatingAnimation bool
}

// returns the next frame in the animation sequence
func (A *Animation) getNextFrame() ([]string, error) {
	if A.CurrentFrame == A.AnimationFrames {
		if !A.RepeatingAnimation {
			return []string{}, AnimationIsOver
		}
		A.CurrentFrame = 0
	}
	A.CurrentFrame++

	return A.Panels[A.CurrentFrame-1], nil
}

type Game struct {
	Config     Config
	State      State
	Charecters Charecters
}

type Charecters struct {
	Char1 Charecter
	Char2 Charecter
}

// used this instead of using p1 , p2 directly in the game struct to make handling inputs easier for functions
// if it ever becomes too bulky to use , feel free to edit
type State struct {
	P1State PlayerState
	P2State PlayerState
}

type AnimationSet map[Action]Animation

// we want a charecter to be an instance of type struct charecter
// we want each charecter to have it's own functions , every charecter have the same set of functions but they are different within
// we want charecters to also hold in values for their own , like a charecter having a mana bar , they must have custom variables
// where would the attack , abilities activation go , is in the client
type Charecter struct {
	Health       int
	Speed        int
	Hitbox       Hitbox
	Attack       func(State)
	Ability1     func(State)
	Entity       map[string]any
	AnimationSet AnimationSet
	ActiveFrame  []string
}

// sprites vs animations
// well for now we don't need sprites , only animations

type Config struct {
	Resulotion      Resulotion `json:"resulotion"`
	MatchDuration   int        `json:"match_duration"`
	SpeedMultiplier int        `json:"speed_multiplier"`
	BackgroundTile  string     `json:"background_tile"` // want to convert to rune , having trouble with that
}

type Resulotion struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type PlayerState struct {
	Position         Point
	Action           Action
	CurrentAnimation *Animation
	Direction        Direction
}

type Point struct {
	X int
	Y int
}

type frame [][]string

// do i have to make the constructer function also add the players ?
// welp you can't have a game without players , none of the methods would work then
// we should throw in the default states , get the chareacters from the input and the config file path too
func NewGame(cfgFilePath string, charecters Charecters) (Game, error) {

	data, err := os.ReadFile(cfgFilePath)

	if err != nil {
		return Game{}, err
	}

	var c Config
	err = json.Unmarshal(data, &c)
	if err != nil {
		return Game{}, err
	}

	return Game{
		Config:     c,
		State:      DefaultState,
		Charecters: charecters,
	}, nil
}

// state should tell the renderer everything it needs to render the scene

// func (game Game) GenerateNextState(gs gameState, inputs playersInputs) playersState  {
// }

// so imagine we have a 3by3 sprite
// let's just see our old code from bombascii

func (f frame) renderSprite(s []string, color string, renderPoint Point) frame {
	outputFrame := f
	midRowIndex := len(s) / 2
	for rowNumber, row := range s {
		relativeY := rowNumber - midRowIndex
		midCharecterIndex := len(row) / 2
		finalY := relativeY + renderPoint.Y
		if (finalY < 0) || (finalY > len(f)-1) {
			continue
		}
		for charecterIndex, charecter := range row {
			relativeX := charecterIndex - midCharecterIndex
			finalX := relativeX + renderPoint.X
			if (finalX > len(f[0])-1) || (finalX < 0) {
				continue
			}
			outputFrame[finalY][finalX] = color + string(charecter) + reset

		}
	}

	return outputFrame
}

// we have a bucket of events we want to use them in order to pretty much alter the state every single time
// so we need something that can construct this events list
func (game Game) MakeFrame(s State) frame {
	// background filling first layer
	// make this into a function btw , generateBasePlate()
	cfg := game.Config    // abbreviation
	res := cfg.Resulotion // abbreviation
	var f frame

	for range res.Height {
		var rowToAdd []string
		for range res.Width {
			rowToAdd = append(rowToAdd, string(cfg.BackgroundTile))
		}
		f = append(f, rowToAdd)
	}
	// second layer for now : rendering the charecters
	position := game.State.P1State.Position
	animationFrame, err := s.P1State.CurrentAnimation.getNextFrame()
	if err != nil {
		fmt.Println(err) // improve the error handling here
	}
	f = f.renderSprite(animationFrame, red, position)
	// next step would be to render the actual charecter instead of this

	return f

}

// optimize this later along with the frame gimmick
func (f frame) RenderFrame() {
	frameToRender := ClearAscii
	for _, row := range f {
		frameToRender += strings.Join(row, "") + "\n"
	}
	fmt.Println(frameToRender)
}
