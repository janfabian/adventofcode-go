package main

import (
	"adventofcode/lib"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	for _, i := range []int{1, 2, 3, 4, 5} {
		wg.Add(1)
		lib.AddTask(func() {
			defer wg.Done()
			fmt.Println("Hello, World!", i)
		})
	}

	wg.Wait()
	lib.Close()
}
