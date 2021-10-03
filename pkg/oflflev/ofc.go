package oflflev

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

// Deuce : 0
const (
	A        = 12
	K        = 11
	Q        = 10
	J        = 9
	T        = 8
	ROYAL    = 9
	STFL     = 8
	QUADS    = 7
	FULL     = 6
	FLUSH    = 5
	STRAIGHT = 4
	TRIPS    = 3
	TWO      = 2
	PAIR     = 1
	HI       = 0
)

type Handrank []int
type Board []Cards

func NewBoard(t, m, b []int) Board {
	return Board{ToCards(t), ToCards(m), ToCards(b)}
}

// ３枚役の強さを評価
// return : [役, ランク, キッカー]
// 数字が大きいほど強い
func EvalTop(cards []int) []int {
	// copy
	cards = append([]int{}, cards...)
	for i := range cards {
		cards[i] %= 13
	}

	sort.Slice(cards, func(i, j int) bool {
		return cards[i] > cards[j]
	})

	if cards[0] == cards[1] && cards[1] == cards[2] {
		return []int{TRIPS, cards[0]} // trips
	} else if cards[0] == cards[1] {
		return []int{PAIR, cards[1], cards[2]} // pair
	} else if cards[1] == cards[2] {
		return []int{PAIR, cards[1], cards[0]} // pair
	} else {
		return []int{HI, cards[0], cards[1], cards[2]}
	}
}

func EvalFive(cards []int) []int {
	// copy
	cards = append([]int{}, cards...)

	flash := isSameSuit(cards)

	for i := range cards {
		cards[i] %= 13
	}

	sort.Slice(cards, func(i, j int) bool {
		return cards[i] > cards[j]
	})

	if cards[0] == A && isStraight(cards) && flash {
		return []int{ROYAL} // Royal flush
	} else if isStraight(cards) && flash {
		return []int{STFL, cards[0]} // straight flush
	} else if cards[0] == A && cards[1] == 3 && cards[2] == 2 && cards[3] == 1 && cards[4] == 0 && flash {
		return []int{STFL, 3} // wheel straight flush
	} else if cards[0] == cards[1] && cards[1] == cards[2] && cards[2] == cards[3] {
		return []int{QUADS, cards[0], cards[4]} // quads
	} else if cards[1] == cards[2] && cards[2] == cards[3] && cards[3] == cards[4] {
		return []int{QUADS, cards[1], cards[0]} // quads
	} else if cards[0] == cards[1] && cards[1] == cards[2] && cards[3] == cards[4] {
		return []int{FULL, cards[0], cards[3]} // full
	} else if cards[0] == cards[1] && cards[2] == cards[3] && cards[3] == cards[4] {
		return []int{FULL, cards[2], cards[0]} // full
	} else if flash {
		return []int{FLUSH, cards[0], cards[1], cards[2], cards[3], cards[4]} // flush
	} else if isStraight(cards) {
		return []int{STRAIGHT, cards[0]} // straight
	} else if cards[0] == 14 && cards[1] == 5 && cards[2] == 4 && cards[3] == 3 && cards[4] == 2 {
		return []int{STRAIGHT, 3} // wheel straight
	} else if cards[0] == cards[1] && cards[1] == cards[2] {
		return []int{TRIPS, cards[0], cards[3], cards[4]} // trips
	} else if cards[1] == cards[2] && cards[2] == cards[3] {
		return []int{TRIPS, cards[1], cards[0], cards[4]} // trips
	} else if cards[2] == cards[3] && cards[3] == cards[4] {
		return []int{TRIPS, cards[2], cards[0], cards[1]} // trips
	} else if cards[0] == cards[1] && cards[2] == cards[3] {
		return []int{TWO, cards[0], cards[2], cards[4]} // twopair
	} else if cards[0] == cards[1] && cards[3] == cards[4] {
		return []int{TWO, cards[0], cards[3], cards[2]} // twopair
	} else if cards[1] == cards[2] && cards[3] == cards[4] {
		return []int{TWO, cards[1], cards[3], cards[0]} // twopair
	} else if cards[0] == cards[1] {
		return []int{PAIR, cards[0], cards[2], cards[3], cards[4]} // pair
	} else if cards[1] == cards[2] {
		return []int{PAIR, cards[1], cards[0], cards[3], cards[4]} // pair
	} else if cards[2] == cards[3] {
		return []int{PAIR, cards[2], cards[0], cards[1], cards[4]} // pair
	} else if cards[3] == cards[4] {
		return []int{PAIR, cards[3], cards[0], cards[1], cards[2]} // pair
	} else {
		return []int{HI, cards[0], cards[1], cards[2], cards[3], cards[4]} // high card
	}
}

func isSameSuit(cards []int) bool {
	suits := []int{}
	for _, card := range cards {
		suits = append(suits, card/13)
	}
	for i := 0; i < len(suits)-1; i++ {
		if suits[i] != suits[i+1] {
			return false
		}
	}
	return true
}

// not wheel
func isStraight(cards []int) bool {
	ranks := []int{}
	for _, card := range cards {
		ranks = append(ranks, card%13)
	}
	sort.Ints(ranks)
	for i := 0; i < len(ranks)-1; i++ {
		if ranks[i+1]-ranks[i] != 1 {
			return false
		}
	}
	return true
}

func CardsToRanks(cards []int) []int {
	ranks := []int{}
	for _, v := range cards {
		ranks = append(ranks, v%13)
	}
	return ranks
}

// a < b : true
func Compair(a, b []int) bool {
	for i := 0; i < min(len(a), len(b)); i++ {
		if a[i] > b[i] {
			return false
		}
	}
	if len(a) < len(b) {
		return true
	} else {
		return false
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

var rankTable = []string{"T", "J", "Q", "K", "A"}
var suitTable = []string{"s", "h", "d", "c"}

type Card int

func (c Card) String() string {
	var rankStr string
	if rank := c % 13; rank < T {
		rankStr = strconv.Itoa(int(rank) + 2)
	} else {
		rankStr = rankTable[rank-T]
	}

	suit := suitTable[int(c)/13]
	return rankStr + suit
}

type Cards []Card

func (c Cards) ToInts() []int {
	s := []int{}
	for _, v := range c {
		s = append(s, int(v))
	}
	return s
}

func (c Cards) String() string {
	// rank順に並べる
	sort.Slice(c, func(i, j int) bool {
		return c[i]%13 > c[j]%13
	})
	strs := []string{}
	for _, v := range c {
		strs = append(strs, v.String())
	}
	return fmt.Sprint(strs)
}

func ToCards(s []int) Cards {
	cards := []Card{}
	for _, v := range s {
		cards = append(cards, Card(v))
	}
	return cards
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Draw(n int) []Card {
	selected := make(map[int]bool)
	for counter := 0; counter < n; {
		a := rand.Intn(52)
		if !selected[a] {
			selected[a] = true
			counter++
		}
	}
	keys := make([]Card, 0, len(selected))
	for k := range selected {
		keys = append(keys, Card(k))
	}
	return keys
}
