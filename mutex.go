package main

import (
    "sync"
	"time"
	"fmt"
)

func main() {
	var mu sync.Mutex
	c := 0
	tStart := time.Now()
	for i := 0; i < 1000000; i++ {
		go func() {
			mu.Lock()
			defer mu.Unlock()
			c++
		}()
	}
	tStop := time.Now()
	time.Sleep(time.Second)
	fmt.Println(c)

	el := tStop.Sub(tStart)
	fmt.Println(el)
}
