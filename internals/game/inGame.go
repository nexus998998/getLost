package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

/*
humble to do list :
1- please order structure of the code to make it more readable

*/

const (
	p1ID = "p1"
	p2ID = "p2"
)

// how do i tell the renderer what's the current frame ?
// simple , currentSprite

type Player int

const (
	Player1 Player = iota
	Player2
)

//go:generate go-enum -f=inGame.go

// ENUM(
// Idle=1
// Walking=2
// Ability1=3
// Ability2=4
// )
type Action int

const (
	_ Action = iota
	Idle
	Walking
	Ability1
	Ability2
	// etc ....
)

//go:generate go-enum -f=inGame.go

// ENUM(
// None
// WalkLeft
// WalkRight
// )
type Input int

const (
	None Input = iota //
	WalkLeft
	WalkRight
	// etc ....
)

type Direction int

const (
	Right Direction = iota
	Left
)

type bindFunc func(p1 PlayerState, p2 PlayerState) (PlayerState, PlayerState)

func (i InputFuncsDependencies) IdleFunc(p1 PlayerState, p2 PlayerState) (PlayerState, PlayerState) {
	p1.Action = Idle
	return p1, p2
}

// we need a function to set the current animation , if it isn't the correct animation it will change it , if it is the animation then it won't
// for that we will pretty much need an identifier to compare them , they are type struct after all

// p1 is the initiator
func (i InputFuncsDependencies) WalkLeftFunc(p1 PlayerState, p2 PlayerState) (PlayerState, PlayerState) {
	p1Pos := &p1.Position
	p1Pos.X -= 1
	if p1Pos.X < 0 {
		p1Pos.X = 0
	}
	p1.Action = Walking
	return p1, p2
}

//
// so what's the solution ? , either make a middle ware that returns a function that takes the state only
// or make a middleware

type binds map[Input]bindFunc

func bindControls(i InputFuncsDependencies) binds {
	// make inputs dependencies
	return map[Input]bindFunc{
		WalkLeft: i.WalkLeftFunc,
		None:     i.IdleFunc,
	}
}

// we are going to attach this to the game , we are going to make a function that attaches global middle ware to all functions

// so we want to feed this function the information that it needs using a
// listening for input

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
			Position: Point{40, 10},
			Action:   Idle,
		},
		P2State: PlayerState{
			Position: Point{30, 10},
			Action:   Idle,
		},
	}

	DefaultCharecter = Charecter{
		Health: 100,
		Speed:  2,
		Hitbox: Hitbox{Point{-1, -1}, Point{1, 1}},
		Name:   "default",
	}

	AnimationIsOver = errors.New("animation is over")
	AssetNotFound   = errors.New("asset not found!")
	TypeMismatchErr = errors.New("expected type doesn't match the asset type ") // this error shall be only used for assets , you can change it later if you used it for anything else

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
type AnimationData struct {
	// Panels             [][]string
	FramesCount        int
	CurrentFrame       int
	RepeatingAnimation bool
}

type panels [][]string

type AnimationConfig struct {
	Panels             [][]string `json:"panels"`
	RepeatingAnimation bool       `json:"repeating_animation"`
	Action             string     `json:"action"`
}

// don't touch this
func (AS *AssetsChunks) LoadAnimationSet(filePath string, ID string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	var cfgs []AnimationConfig
	json.Unmarshal(file, &cfgs)

	for _, cfg := range cfgs {
		AS.loadAnimation(cfg, ID)
	}

	return nil
}

func parseAnimCfg(AC AnimationConfig) (AnimationData, panels) {
	return AnimationData{
		CurrentFrame:       1,
		RepeatingAnimation: AC.RepeatingAnimation,
	}, AC.Panels

}

func FormatFrameKey(ownerID string, action Action, frameNumber int) string {
	return "frame" + action.String() + ownerID + strconv.Itoa(frameNumber)
}

func FormatAnimDataKey(ownerID string, action Action) string {
	return "animData" + action.String() + ownerID
}

// make it go through the animation add the animation and add it's data saperately
func (AS *AssetsChunks) loadAnimation(AC AnimationConfig, ID string) {
	action, err := ParseAction(AC.Action)
	if err != nil {
		panic("this action does not exist , refer to example animation set")
	}

	data, pnls := parseAnimCfg(AC)
	data.FramesCount = len(pnls)
	AS.chunksMap[FormatAnimDataKey(ID, action)] = data
	for frameNum, frame := range pnls {
		AS.chunksMap[FormatFrameKey(ID, action, (frameNum+1))] = frame
	}
}

func (c Charecter) getAnimationSetPath(dataPath string) string {
	formattedName := strings.ReplaceAll(c.Name, " ", "_")
	return fmt.Sprintf("%s/assets/animations/%s.json", dataPath, formattedName)
}

// we need to loop over the list of animations
// and we also need to validate the final map
// maybve we can do this by uhmm a small program inside the assets to validate all the assets

// returns the next frame in the animation sequence

type assetsProvider interface {
	GetFrame(ID string, action Action, frameIndex int) ([]string, error)
	GetAnimData(ID string, action Action) (AnimationData, error)
}

// this assets provider stores needed assets for the match in memory
type AssetsChunks struct {
	chunksMap map[string]any
}

// the idea is to not give direct access to the user for asset chunks ,
func NewAssets() *AssetsChunks {
	return &AssetsChunks{
		chunksMap: map[string]any{},
	}
}

// getAnimation uses the prefix thing for animations
func (AC *AssetsChunks) GetFrame(ID string, action Action, frameIndex int) ([]string, error) {
	anim, ok := AC.chunksMap[FormatFrameKey(ID, action, frameIndex)]
	if !ok {
		return []string{}, AssetNotFound
	}

	if reflect.TypeOf(anim) != reflect.TypeOf([]string{}) {
		return []string{}, TypeMismatchErr
	}

	return anim.([]string), nil
}

func (AC *AssetsChunks) GetAnimData(ID string, action Action) (AnimationData, error) {
	animData, ok := AC.chunksMap[FormatAnimDataKey(ID, action)]
	if !ok {
		return AnimationData{}, AssetNotFound
	}

	if reflect.TypeOf(animData) != reflect.TypeOf(AnimationData{}) {
		return AnimationData{}, TypeMismatchErr
	}

	return animData.(AnimationData), nil
}

type Game struct {
	Config         Config
	State          State
	binds          binds
	AssetsProvider assetsProvider
}

type ActionFunc func(State) (State, error)

// changes the animation to the current action's animation if the action has changed
func handleAnim(AF bindFunc, AP assetsProvider) ActionFunc {

	return func(s State) (State, error) {
		newP1, newP2 := AF(s.P1State, s.P2State)
		if s.P1State.Action != newP1.Action {
			animData, err := AP.GetAnimData(newP1.OwnerID, newP1.Action)
			newP1.AnimState = animData

			if err != nil {
				return State{}, err
			}
		}
		newAnimState, err := newP1.AnimState.progressAnimState()
		newP1.AnimState = newAnimState
		if err != nil {
			return State{}, err
		}
		return State{
			P1State: newP1,
			P2State: newP2,
		}, nil

	}
}

func (g Game) GetBindFunc(I Input, p Player) ActionFunc {

	// so we'd have tjk
	actionFunc := handleAnim(g.binds[I], g.AssetsProvider)

	// so it must be player 2 now
	if p == Player2 {
		actionFunc = func(s State) (State, error) {
			newState, err := actionFunc(s)

			return State{
				P1State: newState.P2State,
				P2State: newState.P1State,
			}, err
		}

	}

	return actionFunc

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

// we want a charecter to be an instance of type struct charecter
// we want each charecter to have it's own functions , every charecter have the same set of functions but they are different within
// we want charecters to also hold in values for their own , like a charecter having a mana bar , they must have custom variables
// where would the attack , abilities activation go , is in the client
type Charecter struct {
	Health      int
	Speed       int
	Hitbox      Hitbox
	Attack      func(State)
	Ability1    func(State)
	Entity      map[string]any
	ActiveFrame []string
	Name        string
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

// this is the best thing i have been doing ever
// have somewhat like an asset for animations , it will load the assets for the charecters at the start of the game , this will also allow for custom assets without recompiling the game
// pretty much an essantial for this
// since we are using a function to load them , we can add some syntax proccessing to it and change it's format to something better usable

type PlayerState struct {
	Position  Point
	Action    Action
	Direction Direction
	AnimState AnimationData
	OwnerID   string
}

type Point struct {
	X int
	Y int
}

type frame [][]string

type InputFuncsDependencies struct {
	Cfg Config
}

func (AS AnimationData) progressAnimState() (AnimationData, error) {
	AS.CurrentFrame++
	if AS.CurrentFrame > AS.FramesCount {
		if !AS.RepeatingAnimation {
			return AnimationData{}, AnimationIsOver
		}
		AS.CurrentFrame = 1
	}

	return AS, nil
}

// do i have to make the constructer function also add the players ?
// welp you can't have a game without players , none of the methods would work then
// we should throw in the default states , get the chareacters from the input and the config file path too
// assets and pres
func NewGame(dataPath string, charecters Charecters) (Game, error) {
	// composing the game\
	cfgFilePath := dataPath + "/config.json"
	data, err := os.ReadFile(cfgFilePath)

	if err != nil {
		return Game{}, err
	}

	var c Config
	err = json.Unmarshal(data, &c)
	if err != nil {
		return Game{}, err
	}

	i := InputFuncsDependencies{
		Cfg: c,
	}

	ap := NewAssets()
	err = ap.LoadAnimationSet(DefaultCharecter.getAnimationSetPath(dataPath), p1ID)
	if err != nil {
		return Game{}, err
	}
	err = ap.LoadAnimationSet(DefaultCharecter.getAnimationSetPath(dataPath), p2ID)
	if err != nil {
		return Game{}, err
	}

	currState := DefaultState
	animS, err := ap.GetAnimData(p1ID, Idle)
	animS2, err := ap.GetAnimData(p2ID, Idle)
	if err != nil {
		return Game{}, err
	}
	currState.P1State.AnimState = animS
	currState.P1State.OwnerID = p1ID
	currState.P2State.AnimState = animS2
	currState.P2State.OwnerID = p2ID
	g := Game{
		Config:         c,
		State:          currState,
		binds:          bindControls(i),
		AssetsProvider: ap,
	}

	return g, nil

}

//
// state should tell the renderer everything it needs to render the scene

func (game Game) GenerateNextState(s State, inputs Input) State {
	NewState, err := game.GetBindFunc(inputs, Player1)(s)
	if err != nil {
		// do logging here
		log.Fatal("error while executing the bind func to generate next state", err.Error())
	}
	return NewState
}

// so imagine we have a 3by3 sprite
// let's just see our old code from bombascii

// this is a function i took from an old ascii project
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
	animationFrame, err := game.AssetsProvider.GetFrame(game.State.P1State.OwnerID, game.State.P1State.Action, game.State.P1State.AnimState.CurrentFrame)
	if err != nil {
		fmt.Println(err) // improve the error handling here
	}
	f = f.renderSprite(animationFrame, green, position)
	// next step would be to render the actual charecter instead of this

	return f

}

// optimize this later along with the frame gimmick
func (f frame) RenderFrame() {
	frameToRender := ClearAscii
	for _, row := range f {
		frameToRender += strings.Join(row, "") + "\r\n"
	}
	fmt.Println(frameToRender)
}
