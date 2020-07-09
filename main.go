package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

var (
	players   = make([]Player, 3)
	bux       = '£'
	buildCost = 100
)

type Player struct {
	id     int
	inPlay bool
	square int
	money  int
	jailed bool
}

func (p *Player) GiveMoney(n int) {
	p.money += n
	if p.money <= 0 {
		fmt.Printf("Player %d has gone bankrupt\n", p.id)
		p.inPlay = false
	}
}

func (p *Player) GoToJail() {
	p.square = squareByName("Jail")
	p.jailed = true
}

func (p *Player) Owned() []*Square {
	var o []*Square
	for i := range board {
		if board[i].owner == p.id {
			o = append(o, &board[i])
		}
	}
	return o
}

func (p *Player) InfoString() string {
	return fmt.Sprintf("%c%d %v", bux, p.money, p.Owned())
}

func roll(n int) int {
	x := 0
	for ; n >= 0; n-- {
		x += 1 + rand.Intn(6)
	}
	return x
}

func main() {
	log.SetFlags(0)
	rand.Seed(time.Now().UnixNano())

	for i := range players {
		players[i].id = i
		players[i].inPlay = true
		players[i].money = 1500
	}

	if false {
		ui()
		return
	}

	for {
		for i := range players {
			p := &players[i]
			if p.id == 0 {
				continue // Banker
			}

			if false {
				if p.id == 1 {
					mstk.Push(mainMenu)
					mstk.Xec(p)
				}
			}

			next := p.square + roll(2)
			if next%len(board) < next {
				log.Printf("player %d passed go\n", p.id)
				p.money += 200
			}
			p.square = next % len(board)
			landed(p, &board[p.square])
			time.Sleep(time.Second / 3)
		}
	}
}
