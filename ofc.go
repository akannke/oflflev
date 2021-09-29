package main

import (
	"reflect"
	"sort"
	"strconv"
	"strings"
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
type Board [][]int

// ３枚役の強さを評価
// return : [役, ランク, キッカー]
// 数字が大きいほど強い
func evalTop(cards []int) []int {
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

func evalFive(cards []int) []int {
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

//    Top Royalties:
//
//    AAA: 22 points
//    KKK: 21 points
//    QQQ: 20 points
//    JJJ: 19 points
//    TTT: 18 points
//    999: 17 points
//    888: 16 points
//    777: 15 points
//    666: 14 points
//    555: 13 points
//    444: 12 points
//    333: 11 points
//    222: 10 points
//    AA: 9 points
//    KK: 8 points
//    QQ: 7 points
//    JJ: 6 points
//    TT: 5 points
//    99: 4 points
//    88: 3 points
//    77: 2 points
//    66: 1 point
func topRoyalty(hr Handrank) int {
	if hr[0] == TRIPS {
		return 10 + hr[1]
	} else if hr[0] == PAIR && hr[1] >= 6 {
		return hr[1] - 3
	} else {
		return 0
	}
}

//// Middle Royalties:
//
//    7-5-4-3-2 perfect: 8 points
//    7-low: 4 points
//    8-low: 2 points
//    9-low: 1 point
func midRoyalty(cards []int) int {
	ranks := cardsToRanks(cards)
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i] > ranks[j]
	})

	perfect := []int{5, 3, 2, 1, 0}
	if reflect.DeepEqual(ranks, perfect) {
		return 8
	} else if cards[0] == 7 {
		return 4
	} else if cards[0] == 8 {
		return 2
	} else if cards[0] == 9 {
		return 1
	} else {
		return 0
	}
}

var botRoyaltyTable []int = []int{0, 0, 0, 0, 2, 4, 6, 10, 15, 25}

//    Bottom Royalties:
//
//    Royal Flush: 25 points
//    Straight Flush: 15 points
//    Four of a Kind: 10 points
//    Full House: 6 points
//    Flush: 4 points
//    Straight: 2 points
func botRoyalty(hr Handrank) int {
	return botRoyaltyTable[hr[0]]
}

func cardsToRanks(cards []int) []int {
	ranks := []int{}
	for _, v := range cards {
		ranks = append(ranks, v%13)
	}
	return ranks
}

func validate(b Board) bool {
	// mid
	midRank := evalFive(b[1])
	if midRank[0] != 0 || midRank[1] > T {
		return false
	}

	topRank := evalTop(b[0])
	botRank := evalFive(b[2])
	return compair(topRank, botRank)
}

// a < b : true
func compair(a, b []int) bool {
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

var rankTable = []string{"T", "Q", "K", "A"}
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

type Row []Card

func (r Row) String() string {
	rowStr := []string{}
	for _, card := range r {
		rowStr = append(rowStr, card.String())
	}

	return strings.Join(rowStr, " ")
}
