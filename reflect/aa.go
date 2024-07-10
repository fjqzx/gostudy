package main

import (
	"encoding/json"
	"fmt"
)

//----结构体标在json中的运用----

type Aaa struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func main() {
	bb := Aaa{"张鑫", "007"}

	//编码的过程 结构体 ———> json
	jsonStr, err := json.Marshal(bb)

	if err != nil {
		fmt.Println("json marshal error", err)
		return
	}

	fmt.Printf("jsonStr = %s\n", jsonStr)

	//解码的过程 jsonstr ———> 结构体
	cc := Aaa{}
	err = json.Unmarshal(jsonStr, &cc)
	if err != nil {
		fmt.Println("json marshal error", err)
		return
	}

	fmt.Printf("OOO = %v\n", cc)

}
