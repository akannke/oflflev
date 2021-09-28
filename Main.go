package main

import "fmt"

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

func main() {
	solve()
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
