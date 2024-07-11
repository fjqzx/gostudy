package main

//----close(?)关闭channel----
//----range----
import "fmt"

func main() {
	c := make(chan int)

	go func() {
		for i := 0; i < 3; i++ {
			c <- i
		}

		close(c)
	}()

	/*	for {
		if data, ok := <-c; ok {
			fmt.Println(data)
		} else {
			break
		}
	}*/

	//使用range来迭代不断操作channel
	for data := range c {
		fmt.Println(data)
	}

	fmt.Println("OOOOO")
}
