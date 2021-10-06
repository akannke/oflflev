package deuce

import (
	"reflect"
	"sort"

	ofc "github.com/akannke/oflflev/pkg/oflflev"
)

type Board ofc.Board

func (b Board) TopRoyalty() int {
	rank := ofc.EvalTop(b.Top)
	return topRoyalty(rank)
}

func (b Board) MiddleRoyalty() int {
	return midRoyalty(b.Middle)
}

func (b Board) BottomRoyalty() int {
	rank := ofc.EvalTop(b.Bottom)
	return botRoyalty(rank)
}

func (b Board) Validate() bool {
	return validate(ofc.Board(b))
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
func topRoyalty(hr ofc.Handrank) int {
	if hr[0] == ofc.TRIPS {
		return 10 + hr[1]
	} else if hr[0] == ofc.PAIR && hr[1] >= 6 {
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
func midRoyalty(cards ofc.Cards) int {
	ranks := ofc.CardsToRanks(cards)
	// decending sort
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i] > ranks[j]
	})

	perfect := []int{5, 3, 2, 1, 0}
	if reflect.DeepEqual(ranks, perfect) {
		return 8
	} else if ranks[0] == 5 {
		return 4
	} else if ranks[0] == 6 {
		return 2
	} else if ranks[0] == 7 {
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
func botRoyalty(hr ofc.Handrank) int {
	return botRoyaltyTable[hr[0]]
}

func validate(b ofc.Board) bool {
	// mid
	midRank := ofc.EvalFive(b.Middle)
	if midRank[0] != 0 || midRank[1] > ofc.T {
		return false
	}

	topRank := ofc.EvalTop(b.Top)
	botRank := ofc.EvalFive(b.Bottom)
	return ofc.Compair(topRank, botRank)
}
