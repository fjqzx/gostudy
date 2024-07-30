package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var c = make(chan int)
var b = 0
var wg sync.WaitGroup

func walkDir(dir string) error {
	// 使用filepath.Walk遍历目录
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // 如果遍历出错，返回错误
		}

		// 判断文件是否是.go文件
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				count := Bb(path)
				c <- count
			}(path)
		}

		return nil // 如果没有错误，返回nil
	})
}

func Bb(path string) int {
	file, err := os.Open(path) /* 打开所要查询的文件 */
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

	return count
}

func Cc() {
	for count := range c {
		b += count
	}
}

func main() {
	// 指定要遍历的目录
	if len(os.Args) < 2 {
		fmt.Println("请提供要遍历的目录作为命令行参数")
		return
	}

	dir := os.Args[1]

	bb := make(chan bool)
	go func() {
		Cc()
		bb <- true
	}()

	// 调用walkDir函数
	if err := walkDir(dir); err != nil {
		fmt.Println("遍历目录时出错:", err)
	}

	wg.Wait()
	close(c)

	<-bb

	fmt.Println("所有代码行数加起来有：", b)
}
