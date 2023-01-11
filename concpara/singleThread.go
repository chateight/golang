package main

import (
	"fmt"
	"time"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(8)
	c := 2
	odd := []int{}
	tStart := time.Now()
	for i := 2; i <= 10000*1000; i++ {
		func() {
			flag := true // if odd number, stay "true"
			for j := 2; j*j <= c; j++ {
				if c%j == 0 {
					flag = false
					break
				}
			}
			if flag == true {
				odd = append(odd, c)
			}
			c++
		}()
	}
	tStop := time.Now()
	fmt.Println(len(odd))

	el := tStop.Sub(tStart)
	fmt.Println(el)
}
