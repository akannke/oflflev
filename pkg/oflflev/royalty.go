package oflflev

type Scorer interface {
	TopRoyalty() int
	MiddleRoyalty() int
	BottomRoyalty() int
	Validate() bool
}

func CalculateScore(scorer Scorer) (bool, int) {
	if scorer.Validate() {
		return true, scorer.TopRoyalty() + scorer.MiddleRoyalty() + scorer.BottomRoyalty()
	} else {
		return false, 0 // faul
	}
}
