package main

import "fmt"

// ----没带缓冲的channel----
func main() {

	c := make(chan int)

	go func() {
		defer fmt.Println("go 运行结束！")

		fmt.Println("go 运行中!")

		c <- 222
	}()

	fmt.Println("main 正在运行中!")

	m := <-c

	fmt.Println("m = ", m)
	fmt.Println("main 运行结束！")
}
