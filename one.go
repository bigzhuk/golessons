package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	t1 := int32(time.Now().Unix())
	for i, arg := range os.Args[0:] {
		fmt.Printf("%d %s\n", i, arg)
		time.Sleep(time.Second)
	}
	t2 := int32(time.Now().Unix())
	fmt.Println(t2 - t1)
}
