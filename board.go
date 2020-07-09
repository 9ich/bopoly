package main

import (
	"fmt"
	"log"
)

type Square struct {
	name      string
	typ       int
	group     int
	cost      int
	buildCost int
	rent      Rent
	owner     int
	houses    int
	mortgaged bool
}

type Rent [6]int // Indexed by properties: 0 1 2 3 4 hotel

func (sq *Square) Rent() int {
	switch {
	case sq.group == SStation:
		return sq.rent[numOwned(SStation, sq.owner)-1]
	case sq.group == SUtility:
		return sq.rent[numOwned(SUtility, sq.owner)-1] * roll(2)
	default:
		return sq.rent[sq.houses]
	}
}

func numOwned(group int, p int) int {
	n := 0
	for i := range board {
		if board[i].group == group && board[i].owner == p {
			n++
		}
	}
	return n
}

const (
	TMisc = iota
	TProp

	SYellow = iota
	SGreen
	SPink
	SRed
	SBrown
	SBlue
	SNavy
	SOrange
	SStation
	SUtility
)

var board = []Square{
	{name: "GO"},
	{name: "Old Kent Road", typ: TProp, group: SBrown, cost: 60, rent: Rent{2, 10, 30, 90, 160, 250}},
	{name: "Community Chest"},
	{name: "Whitechoppel Road", typ: TProp, group: SBrown, cost: 60, rent: Rent{4, 20, 60, 180, 320, 450}},
	{name: "Income Tox"},
	{name: "Kongs Cross Stotion", typ: TProp, group: SStation, cost: 200, rent: Rent{25, 50, 100, 200}},
	{name: "Chonce"},
	{name: "The Ongel Oslington", typ: TProp, group: SBlue, cost: 100, rent: Rent{6, 30, 90, 270, 400, 550}},
	{name: "Euston Road", typ: TProp, group: SBlue, cost: 100, rent: Rent{6, 30, 90, 270, 400, 550}},
	{name: "Pentonville Road", typ: TProp, group: SBlue, cost: 120, rent: Rent{8, 40, 100, 300, 450, 600}},
	{name: "Joil"},
	{name: "Poll Moll", typ: TProp, group: SPink, cost: 140, rent: Rent{10, 50, 150, 450, 625, 750}},
	{name: "Electric Compony", typ: TProp, group: SUtility, cost: 150, rent: Rent{4, 10}},
	{name: "Whiteholl", typ: TProp, group: SPink, cost: 140, rent: Rent{10, 50, 150, 450, 625, 750}},
	{name: "Northumberlond Ovenue", typ: TProp, group: SPink, cost: 160, rent: Rent{12, 60, 180, 500, 700, 900}},
	{name: "Marybone Stotion", typ: TProp, group: SStation, cost: 200, rent: Rent{25, 50, 100, 200}},
	{name: "Bow Street", typ: TProp, group: SOrange, cost: 180, rent: Rent{14, 70, 200, 550, 750, 950}},
	{name: "Community Chest"},
	{name: "Marlborough Street", typ: TProp, group: SOrange, cost: 180, rent: Rent{14, 70, 200, 550, 750, 950}},
	{name: "Vine Street", typ: TProp, group: SOrange, cost: 200, rent: Rent{16, 80, 220, 600, 800, 1000}},
	{name: "Free Porking"},
	{name: "The Strond", typ: TProp, group: SRed, cost: 220, rent: Rent{18, 90, 250, 700, 875, 1050}},
	{name: "Chonce"},
	{name: "Beep Street", typ: TProp, group: SRed, cost: 220, rent: Rent{18, 90, 250, 700, 875, 1050}},
	{name: "Trofolgar Square", typ: TProp, group: SRed, cost: 240, rent: Rent{20, 100, 300, 750, 925, 1100}},
	{name: "Fenchorch St Stotion", typ: TProp, group: SStation, cost: 200, rent: Rent{25, 50, 100, 200}},
	{name: "Leicester Square", typ: TProp, group: SYellow, cost: 260, rent: Rent{22, 110, 330, 800, 975, 1150}},
	{name: "Coventry Street", typ: TProp, group: SYellow, cost: 260, rent: Rent{22, 110, 330, 800, 975, 1150}},
	{name: "Water Works", typ: TProp, group: SUtility, cost: 150, rent: Rent{4, 10}},
	{name: "Piccodolly", typ: TProp, group: SYellow, cost: 280, rent: Rent{22, 120, 360, 850, 1025, 1200}},
	{name: "Go to Joil"},
	{name: "Regent Street", typ: TProp, group: SGreen, cost: 300, rent: Rent{26, 130, 390, 900, 1100, 1275}},
	{name: "Oinksford Street", typ: TProp, group: SGreen, cost: 300, rent: Rent{26, 130, 390, 900, 1100, 1275}},
	{name: "Community Chest"},
	{name: "Bond Street", typ: TProp, group: SGreen, cost: 320, rent: Rent{28, 150, 450, 1000, 1200, 1400}},
	{name: "Liverpool Street Stotion", typ: TProp, group: SStation, cost: 200, rent: Rent{25, 50, 100, 200}},
	{name: "Chonce"},
	{name: "Pork Lane", typ: TProp, group: SNavy, cost: 350, rent: Rent{35, 175, 500, 1100, 1300, 1500}},
	{name: "Super Tox"},
	{name: "Mayfail", typ: TProp, group: SNavy, cost: 400, rent: Rent{50, 200, 600, 1400, 1700, 2000}},
}

func landed(p *Player, sq *Square) {
	switch sq.typ {
	case TProp:
		if sq.owner == 0 {
			if p.money > sq.cost {
				p.GiveMoney(-sq.cost)
				sq.owner = p.id
				log.Printf("player %d buys %v\n", p.id, sq.name)
			} else {
				log.Printf("player %d cannot afford %v\n", p.id, sq.name)
			}
		} else if sq.owner != p.id {
			r := sq.Rent()
			log.Printf("player %d lands on %s and owes %c%d rent to player %d\n",
				p.id, sq.name, bux, r, sq.owner)
			players[sq.owner].money += r
			p.money -= r
			if p.money <= 0 {
				panic("bankrupt")
			}
		} else {
			log.Printf("player %d lands on their property %s\n", p.id, sq.name)
		}
	case TMisc:
		log.Printf("player %d lands on %s and no one cares\n", p.id, sq.name)
	default:
		log.Println("????????????")
	}
}

func squareByName(name string) int {
	for i := range board {
		if board[i].name == name {
			return i
		}
	}
	panic(fmt.Errorf("squareByName: no square with name '%s'", name))
}
