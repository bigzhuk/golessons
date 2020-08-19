// Есть источник данных, это может (быть база данных, АПИ м пр.) генерирующий последовательность данных
// c определенной частотой "Ч" операций/секунду,
// в данном примере, это функция producer, выдающая последовательно надор натуральных чисел от 1 до limit.
// Есть приемник данных, проводящий сложные манипуляции с входящими данными и к примеру сохраняющий результат в другую базу данных,
// умеющий обрабатывать входящие данные с частотой "Ч"/N операциций в секунду, в данном примере это функция processor, вычисляющая квадраты входящего значения,
// где с помощью паузы выполнения эмулируется длительное выполнений операции.
// Для эффективного выполнения задачи, требуется согласовать источник данных и примник данных, путем параллельной обработки в потребителе,
// с ограничение степени параллелизма обработки в размере concurrencySize.
// Таким образом потребитель при получении данных запускает не более чем concurrencySize обработчиков,
//
//                               processor
//      producer -> consumer ->  processor -> terminator (выводит на экран результат, в наеш случае, суммы квадратов входящих наруальных чисел)
//                               ...
//                               processor
//
// При возникновении ошибки обработки, требуется отменить все последующие расчеты, и вернуть ошибку
package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	limit           = 1000
	concurrencySize = 5
)

func producer(limit int) chan int {
	tasks := make(chan int)
	rand.Seed(time.Now().Unix())
	go func() {
		for {
			task := rand.Intn(limit)
			tasks <- task
			time.Sleep(time.Second)
		}
	}()

	return tasks
}

func processor(i int) (int, error) {
	if i == 10 {
		return 0, errors.New("i hate 10")
	}
	time.Sleep(5 * time.Second)
	return i * i, nil
}

func terminator(results chan int) {
	for res := range results {
		fmt.Println(res)
	}
}

func consumer(tasks chan int, results chan int) {
	var procCount int32 = 0
	for {
		num := <-tasks
		go func() {
			if atomic.LoadInt32(&procCount) == concurrencySize {
				tasks <- num
				return
			}
			atomic.AddInt32(&procCount, 1)
			defer atomic.AddInt32(&procCount, -1)
			res, err := processor(num)
			if err != nil {
				panic(err)
			}
			results <- res
		}()
	}
}

func main() {
	tasks := producer(limit)
	res := make(chan int)
	go func() {
		consumer(tasks, res)
	}()

	terminator(res)
}
