package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	file, err := os.Create("hello.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch) // Запуск go-подпрограммы
	}
	for range os.Args[1:] {
		//fmt.Println(<-ch) // Получение из канала ch
		file.WriteString(<-ch)

	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	if strings.HasPrefix(url, "https://") != true && strings.HasPrefix(url, "http://") != true {
		url = "http://" + url
	}
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // Отправка в канал ch return
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // Исключение утечки ресурсов
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
