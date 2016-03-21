package main

import (
	"fmt"
)

type Animal interface {
	Say()
}

type Dog struct {
	name string
}

func (d *Dog) Say() {
	fmt.Println("i'm a dog.....%s", d.name)
}

func main() {
	animals := make(map[string]Animal)
	animals["xiaohua"] = &Dog{"xiaohua"}
	animals["xiaojiji"] = &Dog{"xiaojiji"}
	animals["xiaoheihei"] = &Dog{"xiaoheihei"}
	// animals := []Animal{&Dog{"xiaoming"}, &Dog{"xiaohua"}}
	for _, value := range animals {
		value.Say()
	}
}
