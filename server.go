package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Field
	W = 40
	H = 20

	// Timer / game speed
	TickMs = 80 // milliseconds between ticks

	// Invaders grid
	InvRows   = 3 // number of rows of invaders
	InvStartX = 4 // left offset when spawning invaders
	InvStepX  = 3 // step by X (horizontal spacing)
	InvStartY = 1 // first row Y
	InvEndX   = W - 4

	// Invaders movement
	InvMoveEvery      = 10 // every N-th tick invaders move
	InvLeftBound  int = 1  // left boundary for invader movement
	InvRightBound int = W - 2

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

type Input struct {
	Left  bool `json:"left"`
	Right bool `json:"right"`
	Shoot bool `json:"shoot"`
}

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

type Invader struct {
	X, Y  int
	Alive bool
}
type Bullet struct {
	X, Y, Dy   int
	FromPlayer bool
}

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

func NewGame() *Game {
	g := &Game{
		playerX: PlayerStartX,
		playerY: PlayerStartY,
		lives:   PlayerStartLives,
		invDir:  1,
	}

	// spawn invaders grid, используем константы InvRows, InvStartX, InvStepX
	for y := InvStartY; y <= InvStartY+InvRows-1; y++ {
		for x := InvStartX; x < InvEndX; x += InvStepX {
			g.invaders = append(g.invaders, Invader{X: x, Y: y, Alive: true})
		}
	}
	return g
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // dev only
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer conn.Close()

	game := NewGame()
	inputCh := make(chan Input, 16)

	// reader goroutine: получает входные JSON от клиента
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				// client closed or error
				close(inputCh)
				return
			}
			var in Input
			if err := json.Unmarshal(msg, &in); err != nil {
				continue
			}
			// non-blocking send (drop if buffer full)
			select {
			case inputCh <- in:
			default:
			}
		}
	}()

	ticker := time.NewTicker(time.Duration(TickMs) * time.Millisecond)
	defer ticker.Stop()

	var currInput Input
	for range ticker.C {
		// drain inputCh to get the latest input state
	DrainLoop:
		for {
			select {
			case in, ok := <-inputCh:
				if !ok {
					// connection closed
					return
				}
				currInput = in
			default:
				break DrainLoop
			}
		}

		// apply input
		if currInput.Left && game.playerX > 0 {
			game.playerX--
		}
		if currInput.Right && game.playerX < W-1 {
			game.playerX++
		}
		if currInput.Shoot && game.playerCool == 0 {
			game.bullets = append(game.bullets, Bullet{X: game.playerX, Y: game.playerY - 1, Dy: PlayerBulletDy, FromPlayer: true})
			game.playerCool = PlayerCoolMax
		}
		if game.playerCool > 0 {
			game.playerCool--
		}

		game.tick++

		// invaders group move occasionally
		if game.tick%InvMoveEvery == 0 {
			edge := false
			for i := range game.invaders {
				if !game.invaders[i].Alive {
					continue
				}
				nx := game.invaders[i].X + game.invDir
				if nx < InvLeftBound || nx > InvRightBound {
					edge = true
					break
				}
			}
			if edge {
				for i := range game.invaders {
					if game.invaders[i].Alive {
						game.invaders[i].Y++
					}
				}
				game.invDir *= -1
			} else {
				for i := range game.invaders {
					if game.invaders[i].Alive {
						game.invaders[i].X += game.invDir
					}
				}
			}
		}

		// update bullets
		newBul := make([]Bullet, 0, len(game.bullets))
		for _, b := range game.bullets {
			b.Y += b.Dy
			if b.Y < 0 || b.Y >= H {
				continue
			}
			hit := false
			if b.FromPlayer {
				for i := range game.invaders {
					iv := &game.invaders[i]
					if iv.Alive && iv.X == b.X && iv.Y == b.Y {
						iv.Alive = false
						game.score += ScorePerInvader
						hit = true
						break
					}
				}
			} else {
				// invader bullet can hit player
				if b.X == game.playerX && b.Y == game.playerY {
					game.lives--
					hit = true
				}
			}
			if !hit {
				newBul = append(newBul, b)
			}
		}
		game.bullets = newBul

		// invaders shoot randomly (шанс задаётся через константу)
		if rand.Intn(100) < InvShootChancePercent {
			alive := []Invader{}
			for _, iv := range game.invaders {
				if iv.Alive {
					alive = append(alive, iv)
				}
			}
			if len(alive) > 0 {
				iv := alive[rand.Intn(len(alive))]
				game.bullets = append(game.bullets, Bullet{X: iv.X, Y: iv.Y + 1, Dy: InvaderBulletDy, FromPlayer: false})
			}
		}

		// prepare GameState to send
		st := GameState{
			Width:   W,
			Height:  H,
			PlayerX: game.playerX,
			PlayerY: game.playerY,
			Score:   game.score,
			Lives:   game.lives,
		}
		for _, iv := range game.invaders {
			if iv.Alive {
				st.Invaders = append(st.Invaders, [2]int{iv.X, iv.Y})
			}
		}
		for _, b := range game.bullets {
			st.Bullets = append(st.Bullets, [2]int{b.X, b.Y})
		}

		// send state
		if err := conn.WriteJSON(st); err != nil {
			return
		}

		// simple end condition
		if game.lives <= 0 {
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"gameOver":true}`))
			return
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/ws", wsHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
