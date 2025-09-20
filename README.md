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
	// Поле
	W = 40
	H = 20

	// Таймер / скорость
	TickMs = 80 // миллисекунд между тиками

	// Сетка захватчиков
	InvRows    = 3  // количество рядов захватчиков
	InvStartX  = 4  // левый отступ при создании захватчиков
	InvStepX   = 3  // шаг по X (горизонтальная плотность)
	InvStartY  = 1  // верхняя строка Y
	InvEndX    = W - 4

	// Движение захватчиков
	InvMoveEvery          = 10 // каждый N-ый тик группа захватчиков двигается
	InvLeftBound  int     = 1  // левая граница для движения захватчиков
	InvRightBound int     = W - 2

	// Стрельба захватчиков
	InvShootChancePercent = 10 // шанс в процентах, что кто-то из живых захватчиков выстрелит в тик

	// Игрок
	PlayerStartLives = 3
	PlayerStartX     = W / 2
	PlayerStartY     = H - 2
	PlayerCoolMax    = 6 // тиков между выстрелами игрока

	// Пули
	PlayerBulletDy  = -1
	InvaderBulletDy = 1

	// Очки
	ScorePerInvader = 10
)

```
