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
		// 13 = top + middle + bottom
		for a := range combinations(cards.toInts(), 13, buf) {
			// 8 = top + middle
			for b := range combinations(a, 8, buf) {
				bottom := diff(a, b)
				for top := range combinations(b, 3, buf) {
					mid := diff(b, top)
					c <- NewBoard(top, mid, bottom)
				}
			}
		}
	}()
	return
}

func findBoardTakeBestScore(cards Cards, buf int) (int, Board) {
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
		if loop >= iteration {
			return
		}
		ev = (ev*float64(loop-1) + float64(result.score)) / float64(loop)

		fmt.Println("iteration:", loop)
		fmt.Println("Score:", result.score)
		fmt.Println("Board:", result.board)
		fmt.Println("Received:", result.received)
		fmt.Println("EV:", ev)
		fmt.Println("Elapsed:", time.Since(now))
		fmt.Println("********************")
	}
}

func solve() {
	calcEv(10000)
}

func main() {
	solve()
}
