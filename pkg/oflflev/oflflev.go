package oflflev

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
func Combinations(list []int, select_num, buf int) (c chan []int) {
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
				for sub_comb := range Combinations(list[i+1:], select_num-1, buf) {
					c <- append([]int{list[i]}, sub_comb...)
				}
			}
		}
	}()
	return
}

// 引数に渡すスライスには重複がないこと
func Diff(left, right []int) []int {
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
	Score    int
	Board    Board
	Received Cards
}
