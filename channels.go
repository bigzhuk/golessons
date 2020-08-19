package main

import (
	"fmt"
	"sync"
	"time"
)

func proc(c chan int) {
	for i := 0; i <= 10; i++ {
		c <- i * i
		time.Sleep(time.Second)
	}
	close(c) // close channel

}

func cons(c chan int, wg *sync.WaitGroup) {
	for val := range c {
		time.Sleep(time.Second * 5)
		fmt.Println(val)
	}
	wg.Done()
}

func main() {
	fmt.Println("main() started")
	c := make(chan int)

	var wg sync.WaitGroup // create waitgroup (empty struct)
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go cons(c, &wg)
	}

	proc(c) // start goroutine

	wg.Wait()
	fmt.Println("main() stopped")
}
