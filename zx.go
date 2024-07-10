/*
* @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
* @Date: 2024-07-06 15:13:54
* @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
* @LastEditTime: 2024-07-06 15:17:21
* @FilePath: \zxgo\zx.go
* @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import "fmt"

type baa int

type bans struct {
	bss  string
	auth string
}

func (this bans) Show() {
	fmt.Println("auth = ", this.auth)
	fmt.Println("bss = ", this.bss)
}

func (this bans) Getauth() string {
	return this.auth
}

func (this *bans) Setauth(auth string, bss string) {
	this.auth = auth
	this.bss = bss
}

type Hjm struct {
	bans
	haotin string
}

type Vn struct {
	HH string
}

func (this Vn) GetHH() {
	fmt.Println("HH = ", this.HH)
}

func (this *Hjm) Setauth(auth string, bss string, haotin string) {
	this.auth = auth
	this.bss = bss
	this.haotin = haotin
}

func (this Hjm) Show() {
	fmt.Println("auth = ", this.auth)
	fmt.Println("bss = ", this.bss)
	fmt.Println("haotin = ", this.haotin)
}

func svd(a *int, b *int) {
	var vb int
	vb = *a
	*a = *b
	*b = vb
}

func fg(book1 *bans) {
	book1.bss = "AAA"
	book1.auth = "BBB"
}

func main() {
	//var le int = 2
	//var ba int = 3
	//var a baa = 10

	s := Vn{"zhang"}

	s.GetHH()

	bans := bans{auth: "zhang", bss: "xin"}
	bans.Show()
	bans.Setauth("xiao", "zhang")
	bans.Show()

	var ss Hjm

	ss.auth = "AAA"
	ss.bss = "BBB"
	ss.haotin = "CCC"

	ss.Show()

	//var book1 = bans
	//book1.bss = "dfsf"
	//book1.auth = "zhang"
	//fmt.Printf("%v\n", book1)
	//fg(&book1)

	//svd(&le, &ba)
	//
	//fmt.Printf("%v\n", book1)
	//
	//fmt.Println("a = ", a)
	//
	//fmt.Printf("type of a = %T\n", a)
	//
	//fmt.Println("le = ", le, "ba = ", ba)
}
