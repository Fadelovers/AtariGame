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
