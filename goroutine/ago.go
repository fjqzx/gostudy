package main

//----创建goroutine----
import "time"

func main() {
	go func(a int, b int) bool {
		defer println("a = ", a, "b = ", b)
		return true
	}(5, 6)

	for {
		time.Sleep(1 * time.Second)
	}

}
