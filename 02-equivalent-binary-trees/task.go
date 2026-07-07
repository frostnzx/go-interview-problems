package main

import (
	"golang.org/x/tour/tree"
	"slices"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	ch <- t.Value
	if t.Left != nil {
		go Walk(t.Left, ch)
	}
	if t.Right != nil {
		go Walk(t.Right, ch)
	}
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	go Walk(t1, ch1)
	go Walk(t2, ch2)

	var slice1, slice2 []int
	for {
		v1, _ := <-ch1
		v2, _ := <-ch2
		slice1 = append(slice1, v1)
		slice2 = append(slice2, v2)
		if len(slice1) == 10 && len(slice2) == 10 {
			break
		}
	}
	slices.Sort(slice1)
	slices.Sort(slice2)

	for i := 0; i < 10; i++ {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}
