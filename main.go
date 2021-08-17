/*package main

import (
	"fmt"
)

func intersection(a, b []int) []int {
	result := make([]int, 0, len(a))
	for _, v := range a {
		if IsInSlice(b, v) {
			result = append(result, v)
		}
	}
	return result
}

func main() {

	a := []int{23, 3, 1, 2}
	b := []int{6, 2, 4, 23}
	fmt.Printf("%v\n", intersection(a, b))

	a = []int{1, 1, 1}
	b = []int{1, 1, 1, 1}
	fmt.Printf("%v\n", intersection(a, b))
}
func IsInSlice(a []int, e int) bool {
	for _, v := range a {
		if v == e {
			return true
		}
	}
	return false
}
*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func randNumsGenerator(n int) <-chan int {
	out := make(chan int)
	rand.Seed(time.Now().UnixNano())
	go func() {
		for i := 0; i < n; i++ {
			out <- rand.Intn(1000)
		}
		close(out)
	}()
	return out
}

func main() {
	for num := range randNumsGenerator(10) {
		fmt.Println(num)
	}
}
