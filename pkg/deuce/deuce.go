package deuce

import (
	"fmt"
	"time"

	ofc "github.com/akannke/oflflev/pkg/oflflev"
)

func selectBoardCards(cards ofc.Cards) (c chan ofc.Board) {
	c = make(chan ofc.Board)
	go func() {
		defer close(c)
		for mid := range selectMiddleCards(cards) {
			a := ofc.Diff(cards.ToInts(), mid.ToInts())
			for b := range ofc.Combinations(a, 8, 1) {
				for top := range ofc.Combinations(b, 3, 1) {
					bot := ofc.Diff(b, top)
					c <- ofc.NewBoard(top, mid.ToInts(), bot)
				}
			}
		}
	}()
	return
}

func takeLow(cards ofc.Cards) ofc.Cards {
	lowCards := ofc.Cards{}
	ranks := ofc.CardsToRanks(cards)
	for i := range cards {
		if ranks[i] <= ofc.T {
			lowCards = append(lowCards, cards[i])
		}
	}
	return lowCards
}

func validateMiddle(cards ofc.Cards) bool {
	midRank := ofc.EvalFive(cards)
	if midRank[0] != 0 || midRank[1] > ofc.T {
		return false
	} else {
		return true
	}
}

func selectMiddleCards(cards ofc.Cards) (c chan ofc.Cards) {
	c = make(chan ofc.Cards)
	lowCards := takeLow(cards)
	go func() {
		defer close(c)
		// buf消したい
		for candidate := range ofc.Combinations(lowCards.ToInts(), 5, 1) {
			if mid := ofc.ToCards(candidate); validateMiddle(mid) {
				// 役無しT以下のみ通過
				c <- mid
			}
		}
	}()
	return
}

func findBoardTakeBestScore(cards ofc.Cards) (int, ofc.Board) {
	// todo: bottomをfilterする
	c := selectBoardCards(cards)
	maxScore := 0
	var board ofc.Board = ofc.NewBoard([]int{}, []int{}, []int{})

	for b := range c {
		ok, currentScore := ofc.CalculateScore(Board(b))
		if ok && currentScore >= maxScore {
			maxScore = currentScore
			board = b
		}
	}
	return maxScore, board
}

func findBoardWorker(resultCh chan<- ofc.Result, cardsCh <-chan ofc.Cards) {
	// 適当
	for cards := range cardsCh {
		// 一番良いボード配置を返す関数
		score, board := findBoardTakeBestScore(cards)
		resultCh <- ofc.Result{
			Score:    score,
			Board:    board,
			Received: cards,
		}
	}
}

var now = time.Now()

func CalcEv(iteration int, numDealt int) {
	numWorker := 8
	cardsCh := make(chan ofc.Cards, numWorker)
	resultCh := make(chan ofc.Result)
	defer close(cardsCh)
	defer close(resultCh)

	for i := 0; i < numWorker; i++ {
		go findBoardWorker(resultCh, cardsCh)
	}

	go func() {
		for i := 0; i < iteration; i++ {
			cardsCh <- ofc.Draw(numDealt)
		}
	}()

	var ev float64 = 0
	loop := 0
	for result := range resultCh {
		loop++
		ev = (ev*float64(loop-1) + float64(result.Score)) / float64(loop)

		fmt.Println("iteration:", loop)
		fmt.Println("Score:", result.Score)
		fmt.Println("Board:", result.Board)
		fmt.Println("Received:", result.Received)
		fmt.Println("EV:", ev)
		fmt.Println("Elapsed:", time.Since(now))
		fmt.Println("********************")

		if loop >= iteration {
			return
		}
	}
}

func Solve(iteration int, numDealt int) {
	CalcEv(iteration, numDealt)
}
