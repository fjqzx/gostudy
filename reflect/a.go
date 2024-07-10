package main

import (
	"fmt"
	"reflect"
)

//----结构体标签----

type Aa struct {
	name string `info:"张" doc:"鑫"`
	id   int    `info:"007" doc:"机密"`
}

func findTag(str interface{}) {
	t := reflect.TypeOf(str).Elem()

	for i := 0; i < t.NumField(); i++ {
		taginfo := t.Field(i).Tag.Get("info")
		tagdoc := t.Field(i).Tag.Get("doc")
		fmt.Println("info:", taginfo, "doc", tagdoc)
	}

}

func main() {
	var bb Aa
	findTag(&bb)
}
