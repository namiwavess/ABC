package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"flag"
	"strings"

	melody "gopkg.in/olahol/melody.v1"
)

const (
	foodSize       = 1.0
	lossRate       = 0.002
	absorptionRate = 0.7
	eatingRate     = 0.05

	minPlayerSize = 2.0
	maxPlayerSize = 500.0
)

var maxPlayers int
var worldWidth int
var worldHeight int

type Pos struct{ X, Y int }

type Player struct {
	name    string
	pos     Pos
	mass    float64
	session *melody.Session
}

func detectCollison(p1 Pos, r1 float64, p2 Pos, r2 float64) bool {
	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	d := math.Sqrt(float64(dx*dx + dy*dy))

	if d < r1+r2 {
		return true
	}

	return false
}

func insideCircle(p1 Pos, r1 float64, p2 Pos, r2 float64) bool {

	if r2 < r1 {
		r1, r2 = r2, r1
	}

	dx := p1.X - p2.X
	dy := p1.Y - p2.Y
	d := math.Sqrt(float64(dx*dx + dy*dy))

	if r1-r2 < d {
		return true
	}

	return false
}

func generateFood(n int, w int, h int) []Pos {
	rand.Seed(time.Now().UnixNano())
	food := make([]Pos, n)

	for i := 0; i < n; i++ {
		food[i].X = -1
	}

	return food
}

func addFood(food []Pos, w, h int) map[int]Pos {
	diff := make(map[int]Pos)

	for i, f := range food {
		if f.X == -1 {
			food[i].X = rand.Intn(w)
			food[i].Y = rand.Intn(h)
			diff[i] = food[i]
		}
	}

	return diff
}

func encodeData(diff map[int]Pos, players []Player) string {
	result := new(bytes.Buffer)

	for k, v := range diff {
		fmt.Fprintf(result, "%d,%d,%d,", k, v.X, v.Y)
	}

	fmt.Fprintf(result, " ")

	for i, u := range players {
		if u.pos.X != -1 {
			fmt.Fprintf(result, "%d,%d,%d,%.3f,", i, u.pos.X, u.pos.Y, u.mass)
		}
	}

	return result.String()
}

func (p *Player) eat(o float64) {
	p.mass += o * absorptionRate
	p.mass = math.Min(math.Max(p.mass, minPlayerSize), maxPlayerSize)
}

func (p *Player) loss() {
	p.mass *= (1 - lossRate)
	p.mass = math.Max(p.mass, minPlayerSize)
}

func update(id int, players []Player, food []Pos) {
	pos, r := players[id].pos, players[id].mass

	for i, f := range food {
		if f.X != -1 && detectCollison(pos, r, f, foodSize) {
			food[i].X = -1
			players[id].eat(foodSize)
		}
	}

	for i, u := range players {
		if u.pos.X != -1 && players[i] != players[id] && detectCollison(pos, r, u.pos, u.mass) {
			if r*(1+eatingRate) < u.mass {
				players[i].eat(r)
				players[id].pos.X = -1
				players[id].session.Close()
				fmt.Printf("%s is Dead\n", players[id].name)
			} else if u.mass*1.05 < r {
				players[id].eat(u.mass)
				players[i].pos.X = -1
				players[i].session.Close()
				fmt.Printf("%s is Dead\n", players[i].name)
			}
		}
	}
}

func rangeCheck(p Pos, dx, dy, w, h int) bool {

	x, y := p.X+dx, p.Y+dy

	if x >= 0 && x <= w && y >= 0 && y <= h {
		return true
	}

	return false
}

func initFoodPos(food []Pos) string {
	result := new(bytes.Buffer)

	for i, f := range food {
		if f.X != -1 {
			fmt.Fprintf(result, "%d,%d,%d,", i, f.X, f.Y)
		}
	}

	return result.String()
}

type Act struct {
	i int
	x int
	y int
}

func main() {
	flag.IntVar(&maxPlayers, "num", 100, "set maxPlayers")
	flag.IntVar(&worldWidth, "w", 1000, "set maxPlayers")
	flag.IntVar(&worldHeight, "h", 1000, "set maxPlayers")
	flag.Parse()

	players := make([]Player, maxPlayers)
	queue := make(chan Act, maxPlayers)

	food := generateFood(100, worldWidth, worldHeight)
	addFood(food, worldWidth, worldHeight)
	initFood := initFoodPos(food)

	playerNum := 0
	frame := 0

	waitPlayers := sync.WaitGroup{}
	waitPlayers.Add(maxPlayers)
	mu := sync.Mutex{}

	ws := melody.New()
	view := melody.New()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.HandleRequest(w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("client/dist")))

	http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		view.HandleRequest(w, r)
	})

	view.HandleConnect(func(s *melody.Session) {
		result := new(bytes.Buffer)
		fmt.Fprintf(result, "4 ")

		for i, f := range food {
			if f.X != -1 {
				fmt.Fprintf(result, "%d,%d,%d,", i, f.X, f.Y)
			}
		}

		fmt.Fprintf(result, " ")

		s.Write([]byte(result.String()))
	})

	ws.HandleConnect(func(s *melody.Session) {
		mu.Lock()
		defer mu.Unlock()

		if playerNum < maxPlayers {

			players[playerNum].name = "anonymous"
			players[playerNum].session = s
			players[playerNum].mass = minPlayerSize

			result := new(bytes.Buffer)

			fmt.Fprintf(result, "1 "+strconv.Itoa(playerNum)+" "+initFood+" "+strconv.Itoa(worldWidth)+" "+strconv.Itoa(worldHeight))

			s.Write([]byte(result.String()))
			playerNum++
		} else {
			s.Close()
		}

	})

	ws.HandleMessage(func(s *melody.Session, msg []byte) {
		cmd := strings.Fields(string(msg))

		switch cmd[0] {
		case "2": // Set playerName
			i, _ := strconv.Atoi(cmd[1])
			x, _ := strconv.Atoi(cmd[3])
			y, _ := strconv.Atoi(cmd[4])
			players[i].name = cmd[2]

			players[i].pos.X = x
			players[i].pos.Y = y

			waitPlayers.Done()

		case "3": // Set dx, dy
			i, _ := strconv.Atoi(cmd[1])
			dx, _ := strconv.Atoi(cmd[2])
			dy, _ := strconv.Atoi(cmd[3])

			queue <- Act{i, dx, dy}

		default:
			fmt.Println("unexpected value")
		}

	})

	sendMsg := func() {
		msg := encodeData(addFood(food, worldWidth, worldHeight), players)
		ws.Broadcast([]byte("4 " + msg))
		view.Broadcast([]byte("4 " + msg))
	}

	gameLoop := func() {
		survivor := 0
		frame++

		if len(queue) == 0 {
			sendMsg()
			return
		}

		for j := 0; j < len(queue); j++ {
			act := <-queue
			i, dx, dy := act.i, act.x, act.y

			if players[i].pos.X != -1 && dx*dx+dy*dy <= 9 && rangeCheck(players[i].pos, dx, dy, worldWidth, worldHeight) {
				players[i].pos.X += dx
				players[i].pos.Y += dy

				update(i, players, food)
			}
		}

		for i := 0; i < maxPlayers; i++ {
			if players[i].pos.X != -1 {
				players[i].loss()
				survivor++
			}
		}

		if survivor == 1 {
			for _, u := range players {
				if u.pos.X != -1 {
					fmt.Printf("%s win", u.name)
					os.Exit(0)
				}
			}
		}
		sendMsg()
	}

	go func() {
		waitPlayers.Wait()
		for {
			gameLoop()
			time.Sleep(40 * time.Millisecond)
		}
	}()

	http.ListenAndServe(":3000", nil)
}
