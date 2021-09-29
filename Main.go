package main

import (
	"fmt"
	"github.com/chehsunliu/poker"
)

// Bottom Royalties:

//    Royal Flush: 25 points
//    Straight Flush: 15 points
//    Four of a Kind: 10 points
//    Full House: 6 points
//    Flush: 4 points
//    Straight: 2 points
//
//// Middle Royalties:
//
//    7-5-4-3-2 perfect: 8 points
//    7-low: 4 points
//    8-low: 2 points
//    9-low: 1 point
//
//// Top Royalties:
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
//
// 2: 0, 3: 1, ... , T: 8, J: 9, Q: 10, K: 11, A: 12
// Spade: 0, Heart: 1, Diamond: 2, Clover: 3

func isTrips(cards []int) bool {
	m := make(map[int]int)
	for card := range cards {
		m[card]++
	}

	for _, count := range m {
		if count == 3 {
			return true
		}
	}
	return false
}

func combinations(list []int, select_num, buf int) (c chan []int) {
	c = make(chan []int, buf)
	go func() {
		defer close(c)
		switch {
		case select_num == 0:
			c <- []int{}
		case select_num == len(list):
			c <- list
		case len(list) < select_num:
			return
		default:
			for i := 0; i < len(list); i++ {
				for sub_comb := range combinations(list[i+1:], select_num-1, buf) {
					c <- append([]int{list[i]}, sub_comb...)
				}
			}
		}
	}()
	return
}

func selectBoardCards(cards []int, buf int) (c chan [][]int) {
	c = make(chan [][]int, buf)
	go func() {
		defer close(c)
		// 13 = top + middle + bottom
		for a := range combinations(cards, 13, buf) {
			// 8 = top + middle
			for b := range combinations(a, 8, buf) {
				bottom := diff(a, b)
				for top := range combinations(b, 3, buf) {
					mid := diff(b, top)
					c <- [][]int{top, mid, bottom}
				}
			}
		}
	}()
	return
}

// 引数に渡すスライスには重複がないこと
func diff(left, right []int) []int {
	m := make(map[int]int)
	for _, l := range left {
		m[l]++
	}
	for _, r := range right {
		m[r]--
	}

	result := []int{}
	for i, count := range m {
		if count > 0 {
			result = append(result, i)
		}
	}
	return result
}

func newCards(ss []string) []poker.Card {
	cards := []poker.Card{}
	for _, s := range ss {
		cards = append(cards, poker.NewCard(s))
	}
	return cards
}

func solve() {
	cards := []int{}
	buf := 5
	n := 15
	for i := 0; i < n; i++ {
		cards = append(cards, i)
	}

	c := selectBoardCards(cards, buf)
	b1 := <-c
	b2 := <-c
	fmt.Println(b1)
	fmt.Println(b2)
	fmt.Println(evalTop([]int{0, 13, 26}))
}

func main() {
	solve()
}
