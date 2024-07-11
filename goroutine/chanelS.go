package main

//----带缓冲的channel----
import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int, 3)
	fmt.Println("len(c)=", len(c), "cap(c)=", cap(c))

	go func() {
		defer fmt.Println("go 运行结束")

		for i := 0; i < 4; i++ {
			c <- i
			fmt.Println("go 正在运行， 传入c的值为-->", i, "len(c)=", len(c), "cap(c)=", cap(c))
		}
	}()

	time.Sleep(2 * time.Second)

	for i := 0; i < 4; i++ {
		m := <-c
		fmt.Println("m = ", m)
	}

	fmt.Println("main运行结束")
	time.Sleep(2 * time.Second)
}
