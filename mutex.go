package main

import (
    "sync"
	"time"
	"fmt"
)

//
// to use channel instead of "WaitGroup"
//
// in this case, slower than WaitGroup. Channel may be useful when it takes long processing time and less inforamtion size
//
func main() {
	maxNumber := 1000*10000
	var mu sync.Mutex
	ch := make(chan int)
	defer close(ch)
	c := 2
	oddCh := []int{}
	tStart := time.Now()
	tStop := time.Now()
	for i := 2; i <= maxNumber; i++ {
		go func() {
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
				ch <- c
			}			
			c++
		}()
	}
	for i:= 0; i < maxNumber; i++{
		setBreak := false
		select {
		case p := <- ch:
			oddCh = append(oddCh, p)
			tStop = time.Now()
		case <- time.After(time.Second):		// to check last "ch" data
			fmt.Println("Time Out!")
			setBreak = true
		}
		if setBreak == true{
			break
		}
	}
	
	fmt.Println("len : ", len(oddCh))

	el := tStop.Sub(tStart)
	fmt.Println(el)
}
