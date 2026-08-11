package game

type binds map[Input]bindFunc
type bindFunc func(p1 PlayerState, p2 PlayerState) (PlayerState, PlayerState)

//go:generate go-enum -f=inGame.go

// ENUM(
// None
// WalkLeft
// WalkRight
// )
type Input int

const (
	None Input = iota
	WalkLeft
	WalkRight
	// etc ....
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

type EventCatagory int

const (
	_ EventCatagory = iota
	Movement
	ActivePlayerAction
)

type Player int

const (
	Player1 Player = iota
	Player2
)

type Direction int

const (
	Right Direction = iota
	Left
)

const (
	p1ID = "p1"
	p2ID = "p2"
)

type Hitbox struct {
	UpperLeft   Point
	BottomRight Point
}

type AnimationData struct {
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

type Game struct {
	Config         Config
	State          State
	AssetsProvider assetsProvider
}

type ActionFunc func(State) (State, error)

type Charecters struct {
	Char1 Charecter
	Char2 Charecter
}

type State struct {
	P1State PlayerState
	P2State PlayerState
	Timers  []Timer
}

type event struct {
	FuncID        string
	Catagory      EventCatagory
	PriorityValue int
}

// starts from setoff frame , then when it hits 0 it will fire the event
type Timer struct {
	SetoffFrames int
	CurrFrame    int
	Repeatable   bool
	Event        event
}

type Charecter struct {
	Health int    `json:"heatlh"`
	Speed  int    `json:"speed"`
	Name   string `json:"name"`
}

type MetaData struct {
	AvailableCharecters []string `json:"available_charecters"`
	AvailableStages     []string `json:"available_stages"`
}

type Config struct {
	Resulotion      Resulotion `json:"resulotion"`
	MatchDuration   int        `json:"match_duration"`
	SpeedMultiplier int        `json:"speed_multiplier"`
	BackgroundTile  string     `json:"background_tile"`
}

type Resulotion struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

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
