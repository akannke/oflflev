package main

import "sort"

// Deuce : 0
const (
	A = 12
	K = 11
	Q = 10
	J = 9
	T = 8
)

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
		return []int{3, cards[0]} // trips
	} else if cards[0] == cards[1] {
		return []int{1, cards[1], cards[2]} // pair
	} else if cards[1] == cards[2] {
		return []int{1, cards[1], cards[0]} // pair
	} else {
		return []int{0, cards[0], cards[1], cards[2]}
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
		return []int{9} // Royal flush
	} else if isStraight(cards) && flash {
		return []int{8, cards[0]} // straight flush
	} else if cards[0] == A && cards[1] == 3 && cards[2] == 2 && cards[3] == 1 && cards[4] == 0 && flash {
		return []int{8, 3} // wheel straight flush
	} else if cards[0] == cards[1] && cards[1] == cards[2] && cards[2] == cards[3] {
		return []int{7, cards[0], cards[4]} // quads
	} else if cards[1] == cards[2] && cards[2] == cards[3] && cards[3] == cards[4] {
		return []int{7, cards[1], cards[0]} // quads
	} else if cards[0] == cards[1] && cards[1] == cards[2] && cards[3] == cards[4] {
		return []int{6, cards[0], cards[3]} // full
	} else if cards[0] == cards[1] && cards[2] == cards[3] && cards[3] == cards[4] {
		return []int{6, cards[2], cards[0]} // full
	} else if flash {
		return []int{5, cards[0], cards[1], cards[2], cards[3], cards[4]} // flush
	} else if isStraight(cards) {
		return []int{4, cards[0]} // straight
	} else if cards[0] == 14 && cards[1] == 5 && cards[2] == 4 && cards[3] == 3 && cards[4] == 2 {
		return []int{4, 3} // wheel straight
	} else if cards[0] == cards[1] && cards[1] == cards[2] {
		return []int{3, cards[0], cards[3], cards[4]} // trips
	} else if cards[1] == cards[2] && cards[2] == cards[3] {
		return []int{3, cards[1], cards[0], cards[4]} // trips
	} else if cards[2] == cards[3] && cards[3] == cards[4] {
		return []int{3, cards[2], cards[0], cards[1]} // trips
	} else if cards[0] == cards[1] && cards[2] == cards[3] {
		return []int{2, cards[0], cards[2], cards[4]} // twopair
	} else if cards[0] == cards[1] && cards[3] == cards[4] {
		return []int{2, cards[0], cards[3], cards[2]} // twopair
	} else if cards[1] == cards[2] && cards[3] == cards[4] {
		return []int{2, cards[1], cards[3], cards[0]} // twopair
	} else if cards[0] == cards[1] {
		return []int{1, cards[0], cards[2], cards[3], cards[4]} // pair
	} else if cards[1] == cards[2] {
		return []int{1, cards[1], cards[0], cards[3], cards[4]} // pair
	} else if cards[2] == cards[3] {
		return []int{1, cards[2], cards[0], cards[1], cards[4]} // pair
	} else if cards[3] == cards[4] {
		return []int{1, cards[3], cards[0], cards[1], cards[2]} // pair
	} else {
		return []int{0, cards[0], cards[1], cards[2], cards[3], cards[4]} // high card
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
