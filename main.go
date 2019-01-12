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
var foodquantity int
var cheat bool
var mannedname string

type Pos struct{ X, Y int }

type Mouse struct{ X,Y int }
var mouse Mouse

type Player struct {
	name    string
	pos     Pos
	mass    float64
	session *melody.Session
}

// プレイヤー同士の衝突検出
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

	fmt.Fprintf(result, " ")
	fmt.Fprintf(result, "%d,%d",mouse.X,mouse.Y)

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

var killflag bool = false
var killlog string = ""
var rank int = -1
var killnum int = 0

func update(id int, players []Player, food []Pos) {
	pos, r := players[id].pos, players[id].mass

	for i, f := range food {
		if f.X != -1 && detectCollison(pos, r, f, foodSize) {
			food[i].X = -1
			players[id].eat(foodSize)
		}
	}
	for i, u := range players {
		if cheat && u.name == mannedname && u.mass < 200 {//チート
			u.mass = 200
		}
		if u.pos.X != -1 && players[i] != players[id] && detectCollison(pos, r, u.pos, u.mass) {
			var killmsg = new(bytes.Buffer)
			survivor := 0
			for j := 0; j < maxPlayers; j++ {
				if players[j].pos.X != -1 {
					survivor++
				}
			}
			if r*(1+eatingRate) < u.mass {
				players[i].eat(r)
				players[id].pos.X = -1
				players[id].session.Close()
				if players[id].name == mannedname {
					rank = survivor
				}else if players[i].name == mannedname {
					killnum++
				}
				fmt.Printf("%s is Dead killed by %s\n", players[id].name,players[i].name)
				fmt.Fprintf(killmsg,"%s,is,Dead,killed,by,%s %d %d", players[id].name,players[i].name,killnum,survivor-1)
				killflag = true
				killlog = killmsg.String()
			} else if u.mass*1.05 < r {
				players[id].eat(u.mass)
				players[i].pos.X = -1
				players[i].session.Close()
				if players[i].name == mannedname {
					rank = survivor
				}else if players[id].name == mannedname {
					killnum++
				}
				fmt.Printf("%s is Dead killed by %s\n", players[i].name,players[id].name)
				fmt.Fprintf(killmsg,"%s,is,Dead,killed,by,%s %d %d", players[i].name,players[id].name,killnum,survivor-1)
				killflag = true
				killlog = killmsg.String()
				
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

var continueflag bool = true

func main() {
	flag.IntVar(&maxPlayers, "num", 100, "set maxPlayers")
	flag.IntVar(&worldWidth, "w", 1000, "set worldWidth")
	flag.IntVar(&worldHeight, "h", 1000, "set worldHeight")
	flag.IntVar(&foodquantity, "food", 100, "set food quantity")
	flag.BoolVar(&cheat, "c", false, "cheat mode always your size 200")
	flag.StringVar(&mannedname, "name", "Manned", "Manned players name")
	flag.Parse()

	mouse.X = worldWidth/2
	mouse.Y = worldHeight/2

	players := make([]Player, maxPlayers)
	queue := make(chan Act, maxPlayers)

	food := generateFood(foodquantity, worldWidth, worldHeight)
	addFood(food, worldWidth, worldHeight)
	initFood := initFoodPos(food)

	playerNum := 0
	frame := 0

	waitPlayers := sync.WaitGroup{} 
	//playerの同期
	waitPlayers.Add(maxPlayers)
	//同期するplayerの数
	mu := sync.Mutex{}

	ws := melody.New() //playerとのwebsocket
	view := melody.New() //webブラウザとのwebsocket

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.HandleRequest(w, r)
	})

	http.Handle("/", http.FileServer(http.Dir("client/dist")))

	http.HandleFunc("/view", func(w http.ResponseWriter, r *http.Request) {
		view.HandleRequest(w, r)
	})

	view.HandleConnect(func(s *melody.Session) {
		result := new(bytes.Buffer)
		setinfo := new(bytes.Buffer)
		fmt.Fprintf(setinfo, "9 %d %d",0,maxPlayers)
		view.Broadcast([]byte(setinfo.String()))
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

	view.HandleMessage(func(s *melody.Session, msg []byte) {
		cmd := strings.Fields(string(msg))

		switch cmd[0] {
		case "5": //マウスの座標受け取り
			
			mousex, _ := strconv.Atoi(cmd[1])
			mousey, _ := strconv.Atoi(cmd[2])

			mouse.X = mousex
			mouse.Y = mousey

		case "7": //Quit Button
			for _ , p := range players {
				p.session.Close()
			}
			os.Exit(0)

		default:
			fmt.Println("unexpected value js")
		}
	})

	//クライアント(Cell.py)からメッセージを受信したときに行う処理
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
		if(killflag){
			view.Broadcast([]byte("6 " + killlog))
			killflag = false
		}
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

		if rank != -1 {
			var losemsg = new(bytes.Buffer)
			fmt.Fprintf(losemsg,"<br><h1>#,%d/%d<br></h1><h1,style=\"color:blue\">まあ、こんな日もあるのさ!次はもう少しツイてますように!</h1>", rank, maxPlayers)
			fmt.Fprintf(losemsg,"<br><h2>ランク,#%d　　キル,%d,プレイヤー　　報酬　なし</h1>",rank,killnum)
			view.Broadcast([]byte("8 " + losemsg.String()))
			rank = -1		
		}

		if survivor == 1 && continueflag {
			for _, u := range players {
				if u.pos.X != -1 {
					fmt.Printf("\n%s win\n", u.name)
					var winmsg = new(bytes.Buffer)
					if(u.name == mannedname){
						fmt.Fprintf(winmsg,"<br><h1,style=\"color:orange\">#1,<span,style=\"color:gray\">/%d</span></h1><h1,style=\"color:red\">%s,win!!<br>勝った!勝った!夕飯はドン勝だ!!</h1>", maxPlayers,mannedname)
						fmt.Fprintf(winmsg,"<br><h2>ランク,#1　　キル,%d,プレイヤー　　報酬　単位</h1>",killnum)
						view.Broadcast([]byte("8 " + winmsg.String()))
						}else{
						//fmt.Fprintf(winmsg,"<br><h1>#,%d<br>%s,win</h1><h1,style=\"color:blue\">まあ、こんな日もあるのさ!次はもう少しツイてますように!</h1>", rank,u.name)
						time.Sleep(5 * time.Millisecond)
						fmt.Fprintf(winmsg,"<h1>%s,win</h1>",u.name)
						view.Broadcast([]byte("6 "+ winmsg.String()))
					}
					continueflag = false
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
