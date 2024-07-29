package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var c = make(chan int)
var b = 0

func walkDir(dir string) error {
	// 使用filepath.Walk遍历目录
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // 如果遍历出错，返回错误
		}

		// 判断是否是目录
		if info.IsDir() {
			//fmt.Println("目录:", path)
		} else {
			//文件
			fmt.Println("文件:", path)
			go Cc()
			Bb(path)
		}
		return nil // 如果没有错误，返回nil
	})
}

func Bb(path string) {
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
	c <- count
	fmt.Println("一共有", count, "行")
}

func Cc() {
	n := <-c
	b = b + n
}

func main() {
	// 指定要遍历的目录
	dir := "D:\\zxmxcx2\\ATM\\src\\images"

	// 调用walkDir函数
	if err := walkDir(dir); err != nil {
		fmt.Println("遍历目录时出错:", err)
	}
	fmt.Println("所有代码行数加起来有：", b)
}
