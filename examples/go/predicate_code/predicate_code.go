package main

import "fmt"

func main() {
	p := newParser("", []byte("%"))
	fmt.Println(p.parse(nil))

	p = newParser("", []byte("-"))
	fmt.Println(p.parse(nil))
}
