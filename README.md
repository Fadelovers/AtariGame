# Space Invaders
[![Language](https://img.shields.io/badge/Language-Golang-blue)](https://go.dev/)
[![Language](https://img.shields.io/badge/Language-Html-orange)](https://www.w3schools.com/html/default.asp)
![Language](https://img.shields.io/badge/Type-Web_Game-black)

## About

Space Invaders is a two-dimensional fixed shooter game in which the player controls a ship with lasers by moving it horizontally across the bottom of the screen and firing at descending aliens. The aim is to defeat five rows of ten aliens that move horizontally back and forth across the screen as they advance towards the bottom of the screen. The player defeats an alien, and earns points, by shooting it with the laser cannon. As more aliens are defeated, the aliens' movement and the game's music both speed up. 

![Image](https://github.com/user-attachments/assets/b6b17595-f701-4875-8c5f-16e40a117288)

## How To Play

compile the project and run it on localhost

Use ← → or A/D to move, Space to shoot. Click on page to focus

## project settings

```go
	const (
	// Field
	W = 40
	H = 20

	// Timer / game speed
	TickMs = 80 // milliseconds between ticks

	// Invaders grid
	InvRows    = 3  // number of rows of invaders
	InvStartX  = 4  // left offset when spawning invaders
	InvStepX   = 3  // step by X (horizontal spacing)
	InvStartY  = 1  // first row Y
	InvEndX    = W - 4

	// Invaders movement
	InvMoveEvery          = 10 // every N-th tick invaders move
	InvLeftBound  int     = 1  // left boundary for invader movement
	InvRightBound int     = W - 2

	// Invaders shooting
	InvShootChancePercent = 10 // chance in percent each tick that one invader shoots

	// Player
	PlayerStartLives = 3
	PlayerStartX     = W / 2
	PlayerStartY     = H - 2
	PlayerCoolMax    = 6 // ticks between player shots

	// Bullets
	PlayerBulletDy  = -1
	InvaderBulletDy = 1

	// Scoring
	ScorePerInvader = 10
)


```

## OOP

OOP is used in this code to implement the web interface.


```go
type Input struct {
    Left  bool `json:"left"`
    Right bool `json:"right"`
    Shoot bool `json:"shoot"`
}
```

Purpose: a DTO (data transfer object) for input from the client — keys/control state. JSON tags allow json.Unmarshal to parse incoming JSON messages directly into this struct. This is convenient and safe: you have an explicit list of allowed inputs.

```go
type GameState struct {
    Width    int      `json:"width"`
    Height   int      `json:"height"`
    PlayerX  int      `json:"playerX"`
    PlayerY  int      `json:"playerY"`
    Invaders [][2]int `json:"invaders"`
    Bullets  [][2]int `json:"bullets"`
    Score    int      `json:"score"`
    Lives    int      `json:"lives"`
}

```

Purpose: the structure you serialize and send to the client (via WebSocket). This is a presentation of the game state, simplified and prepared for JSON. Separating the internal state (Game) and GameState lets you avoid exposing unnecessary details and build only the information the frontend needs.

```go
type Invader struct {
    X, Y  int
    Alive bool
}
```
Purpose: a model of a single invader in the internal game state. It stores coordinates and a alive/dead flag. Having the boolean Alive makes it easy to mark an invader as destroyed without removing the element from the slice (removing/inserting is more expensive and complicates indices).

```go
type Bullet struct {
    X, Y, Dy   int
    FromPlayer bool
}
```

Purpose: a model of a projectile. Dy is the vertical direction/speed, FromPlayer distinguishes player bullets from invader bullets (for example, to decide what they can hit).

```go
type Game struct {
    playerX, playerY int
    playerCool       int
    lives            int
    score            int

    invaders []Invader
    bullets  []Bullet
    invDir   int
    tick     int
}

```

Purpose: the main internal state of the game + logic (in your code the logic is implemented directly inside wsHandler, but Game accumulates state). This is the game model: it contains all data that changes during the game.