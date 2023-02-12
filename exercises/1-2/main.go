package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- int) {
	for i := 0; ; i++ {
		ch <- i
		fmt.Println("Produced: ", i)
		time.Sleep(time.Second)
	}
}

func comsumer(ch <-chan int) {
	for {
		fmt.Println("Consumed", <-ch)
		time.Sleep(time.Second)
	}
}

func main() {
	ch := make(chan int, 10)

	go producer(ch)
	go comsumer(ch)

	// Quit by pressing enter
	var input string
	fmt.Scanln(&input)
}
