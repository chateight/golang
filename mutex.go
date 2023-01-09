package main

import (
    "sync"
	"time"
	"fmt"
	"math"
)

func main() {
	var mu sync.Mutex
	c := 2
	odd := []int{}
	tStart := time.Now()
	for i := 2; i <= 10000*1000; i++ {
		go func() {
			flag := true		// if odd number, stay "true"
			mu.Lock()
			defer mu.Unlock()
			for j :=2; j <= int(math.Sqrt(float64(c))); j++ {
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
	tStop := time.Now()
	time.Sleep(time.Second)
	fmt.Println(len(odd))

	el := tStop.Sub(tStart)
	fmt.Println(el)
}
