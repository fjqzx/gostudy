package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("hh.go")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	count := 0
	for scanner.Scan() {
		log.Println(scanner.Text())
		count++
	}

	fmt.Println("一共有", count, "行")
}
