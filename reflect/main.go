package main

//----反射----

import (
	"fmt"
	"reflect"
)

type Fanshe struct {
	Id   int
	Name string
	Jin  float64
}

func (this Fanshe) Aa() {
	fmt.Printf("%v\n", this)
}

func main() {
	aa := Fanshe{Id: 1, Name: "zx", Jin: 1.23}
	//aa.Aa()
	Bb(&aa)
}

func Bb(hh interface{}) {
	a1 := reflect.TypeOf(hh).Elem()
	fmt.Println("a1 = ", a1.Name())

	a2 := reflect.ValueOf(hh).Elem()
	fmt.Println("a2 = ", a2)

	for i := 0; i < a1.NumField(); i++ {
		field := a1.Field(i)
		value := a2.Field(i).Interface()

		fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
	}

	for i := 0; i < a1.NumMethod(); i++ {
		b := a1.Method(i)
		fmt.Printf("methodName:%s methodType:%v\n", b.Name, b.Type)
	}
}
