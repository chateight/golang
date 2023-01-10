package main

import (
    "sync"
	"time"
	"fmt"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	c := 2
	odd := []int{}
	tStart := time.Now()
	for i := 2; i <= 10000*1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			flag := true		// if odd number, stay "true"
			mu.Lock()
			defer mu.Unlock()
			for j :=2; j*j <= c ; j++ {
				if c%j == 0{
					flag = false
					break
				}
			}
			if flag == true{
				odd = append(odd, c)
			}			
			c++
		}()
	}
	wg.Wait()
	tStop := time.Now()
	fmt.Println(len(odd))

	el := tStop.Sub(tStart)
	fmt.Println(el)
}
