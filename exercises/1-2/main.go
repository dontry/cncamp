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
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		fmt.Println("Consumed: ", <-ch)
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
