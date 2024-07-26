package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("hh.go") /* 打开所要查询的文件 */
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file) /* 创建Scanner */

	count := 0 /* 计数器 */
	
	//遍历
	for scanner.Scan() {
		//log.Println(scanner.Text()) /* 打印每行的消息 */
		count++
	}

	fmt.Println("一共有", count, "行")
}
