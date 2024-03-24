package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	for i := 0; i < 15; i++ {
		duration := time.Duration(1+rand.Intn(10)) * time.Second
		time.Sleep(duration)
		fmt.Println(duration)
	}
}
