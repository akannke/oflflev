package regular

import (
	ofc "github.com/akannke/oflflev/pkg/oflflev"
)

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

var midRoyaltyTable []int = []int{0, 0, 0, 2, 4, 8, 12, 20, 30, 50}

// Middle Royalties
// Royal Flush: 50 points
// Straight Flush: 30 points
// Four of a Kind: 20 points
// Full House: 12 points
// Flush: 8 points
// Straight: 4 points
// Three of a Kind: 2 points
func midRoyalty(hr ofc.Handrank) int {
	return midRoyaltyTable[hr[0]]
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
