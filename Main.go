package main

import (
	"fmt"
	"time"
)

// Bottom Royalties:

//    Royal Flush: 25 points
//    Straight Flush: 15 points
//    Four of a Kind: 10 points
//    Full House: 6 points
//    Flush: 4 points
//    Straight: 2 points
//
//
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

// buf消したい、override?
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

func selectBoardCards(cards Cards, buf int) (c chan Board) {

	c = make(chan Board, buf)
	go func() {
		defer close(c)
		for mid := range selectMiddleCards(cards) {
			a := diff(cards.toInts(), mid.toInts())
			for b := range combinations(a, 8, 1) {
				for top := range combinations(b, 3, 1) {
					bot := diff(b, top)
					c <- NewBoard(top, mid.toInts(), bot)
				}
			}
		}
	}()
	return
}

func takeLow(cards Cards) Cards {
	lowCards := Cards{}
	ranks := cardsToRanks(cards.toInts())
	for i := range cards {
		if ranks[i] <= T {
			lowCards = append(lowCards, cards[i])
		}
	}
	return lowCards
}

func selectMiddleCards(cards Cards) (c chan Cards) {
	c = make(chan Cards)
	lowCards := takeLow(cards)
	go func() {
		defer close(c)
		// buf消したい
		for candidate := range combinations(lowCards.toInts(), 5, 1) {
			if mid := toCards(candidate); validateMiddle(mid) {
				// 役無しT以下のみ通過
				c <- mid
			}
		}
	}()
	return
}

func findBoardTakeBestScore(cards Cards, buf int) (int, Board) {
	// todo: bottomをfilterする
	c := selectBoardCards(cards, buf)
	maxScore := 0
	var board Board = nil

	for b := range c {
		if ok, currentScore := calcScore(b); ok && currentScore >= maxScore {
			maxScore = currentScore
			board = b
		}
	}
	return maxScore, board
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

type Result struct {
	score    int
	board    Board
	received Cards
}

func findBoardWorker(resultCh chan<- Result, cardsCh <-chan Cards) {
	// 適当
	buf := 10
	for cards := range cardsCh {
		score, board := findBoardTakeBestScore(cards, buf)
		resultCh <- Result{
			score:    score,
			board:    board,
			received: cards,
		}
	}
}

var now = time.Now()

func calcEv(iteration int) {
	numWorker := 8
	cardsCh := make(chan Cards, numWorker)
	resultCh := make(chan Result)
	defer close(cardsCh)
	defer close(resultCh)

	for i := 0; i < numWorker; i++ {
		go findBoardWorker(resultCh, cardsCh)
	}

	go func() {
		for i := 0; i < iteration; i++ {
			cardsCh <- draw(14)
		}
	}()

	var ev float64 = 0
	loop := 0
	for result := range resultCh {
		loop++
		ev = (ev*float64(loop-1) + float64(result.score)) / float64(loop)

		fmt.Println("iteration:", loop)
		fmt.Println("Score:", result.score)
		fmt.Println("Board:", result.board)
		fmt.Println("Received:", result.received)
		fmt.Println("EV:", ev)
		fmt.Println("Elapsed:", time.Since(now))
		fmt.Println("********************")

		if loop >= iteration {
			return
		}
	}
}

func solve() {
	calcEv(10000)
}

func main() {
	solve()
}
