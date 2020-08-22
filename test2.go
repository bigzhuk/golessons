package main

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

const (
	limit           = 50
	concurrencySize = 5
)

func producer(limit int) chan int {
	inputStream := make(chan int)
	go func() {
		for i := 1; i <= limit; i++ {
			inputStream <- i
		}
		close(inputStream)
	}()
	return inputStream
}

func square(i int) (int, error) {
	if i == 5 {
		return 0, errors.New("i hate 5")
	}
	time.Sleep(5 * time.Second)
	return i * i, nil
}

func concurrentCalculationUnit(ctx context.Context, i int, wg *sync.WaitGroup, semaChan chan struct{}, errorChan chan error) {
	select {
	case semaChan <- struct{}{}:
		wg.Add(1)
	case <-ctx.Done():
		log.Println("skip calculation for", i, "cancelled")
		return
	}

	defer func() {
		<-semaChan
		wg.Done()
	}()

	square, err := square(i)
	if err != nil {
		log.Println("error calculation for", i, err)
		errorChan <- err
		return
	}
	log.Println("square for", i, "is", square)
}

func main() {
	var wg sync.WaitGroup
	semaChan := make(chan struct{}, concurrencySize)
	errorChan := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for _ = range errorChan {
			cancel()
		}
	}()

	for i := range producer(limit) {
		i := i
		go concurrentCalculationUnit(ctx, i, &wg, semaChan, errorChan)
	}
	wg.Wait()
}
