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
	W = 40
	H = 20
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
		playerX: W / 2,
		playerY: H - 2,
		lives:   3,
		invDir:  1,
	}
	// spawn invaders grid
	for y := 1; y <= 3; y++ {
		for x := 4; x < W-4; x += 3 {
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

	// reader goroutine: receives input JSONs from client
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

	ticker := time.NewTicker(80 * time.Millisecond)
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
			game.bullets = append(game.bullets, Bullet{X: game.playerX, Y: game.playerY - 1, Dy: -1, FromPlayer: true})
			game.playerCool = 6
		}
		if game.playerCool > 0 {
			game.playerCool--
		}

		game.tick++

		// invaders group move occasionally
		if game.tick%10 == 0 {
			edge := false
			for i := range game.invaders {
				if !game.invaders[i].Alive {
					continue
				}
				nx := game.invaders[i].X + game.invDir
				if nx < 1 || nx > W-2 {
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
						game.score += 10
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

		// invaders shoot randomly
		if rand.Intn(100) < 10 {
			alive := []Invader{}
			for _, iv := range game.invaders {
				if iv.Alive {
					alive = append(alive, iv)
				}
			}
			if len(alive) > 0 {
				iv := alive[rand.Intn(len(alive))]
				game.bullets = append(game.bullets, Bullet{X: iv.X, Y: iv.Y + 1, Dy: 1, FromPlayer: false})
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
