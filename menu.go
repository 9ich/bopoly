package main

import (
	"fmt"
	"os"
)

var mstk menuStack

type menuFn func(*Player)
type menuStack []menuFn

func (s *menuStack) Push(f menuFn) {
	*s = append(*s, f)
}

func (s *menuStack) Pop() {
	n := len(*s)
	if n < 1 {
		panic("menuStack: underflow")
	}
	*s = (*s)[:n-1]
}

func (s *menuStack) Xec(p *Player) {
	for len(*s) != 0 {
		(*s)[len(*s)-1](p)
	}
}

func mainMenu(p *Player) {
	c := prompt([]string{"Roll dice", "Build", "Mortgage", "Quit"})
	switch c {
	case 1:
		mstk.Pop()
	case 2:
		mstk.Push(buildMenu)
	case 3:
		mstk.Push(mortMenu)
	case 4:
		os.Exit(0)
	}
}

func buildMenu(p *Player) {
	choices := []string{"<Back>"}
	sqrs := p.Owned()
	for _, sq := range sqrs {
		s := fmt.Sprintf("%s (group %d) (%d houses) %c%d",
			sq.name, sq.group, sq.houses, bux, sq.buildCost)
		choices = append(choices, s)
	}
	c := prompt(choices)
	switch c {
	case 1:
		mstk.Pop()
	default:
		sq := &board[squareByName(sqrs[c-1].name)]
		if p.money >= buildCost {
			sq.houses++
			p.GiveMoney(-buildCost)
			mstk.Pop()
		} else {
			fmt.Println("Can't afford that")
		}
	}
}

func mortMenu(p *Player) {
	choices := []string{"<Back>"}
	sqrs := p.Owned()
	for _, sq := range sqrs {
		s := fmt.Sprintf("%s (group %d) (%d houses) %c%d",
			sq.name, sq.group, sq.houses, bux, sq.cost/2)
		choices = append(choices, s)
	}
	c := prompt(choices)
	switch c {
	case 1:
		mstk.Pop()
	default:
		sq := &board[squareByName(sqrs[c-1].name)]
		sq.mortgaged = true
		p.GiveMoney(sq.cost / 2)
		mstk.Pop()
	}
}

func prompt(choices []string) int {
	fmt.Println("Pick one:")
	for i := range choices {
		fmt.Printf("(%d) %s\n", i+1, choices[i])
	}
	fmt.Println(players[1].InfoString())
	fmt.Print("> ")
	var c int
	fmt.Scan(&c)
	if c < 1 || c >= len(choices)+1 {
		fmt.Println("That's not a valid choice")
		return prompt(choices)
	}
	return c
}
