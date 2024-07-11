package main

//----select----
import (
	"fmt"
)

func Bb(c, p chan int) {
	x, y := 1, 1
	for {
		select {
		case c <- x:
			y = x + y
			x = y
		case <-p:
			fmt.Println("结束")
			return
		}
	}
}

func main() {
	c := make(chan int)
	p := make(chan int)
	go func() {
		for i := 0; i < 3; i++ {
			fmt.Println(<-c)
		}

		p <- 0
	}()
	Bb(c, p)
}
