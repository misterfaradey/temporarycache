package main

import (
	"fmt"

	"github.com/misterfaradey/temporarycache"
)

var (
	mem temporarycache.Mem
)

const (
	size = 10
)

func main() {

	if size > 0 {
		mem = temporarycache.InitMemCache(size)
	} else {
		mem = temporarycache.InitMemCache(0)
	}

	mem.Write(1, "1")
	mem.Write(2, "2")
	mem.Write(3, "3")
	mem.Write(4, "4")
	mem.Write(5, "5")
	mem.Write(6, "6")
	mem.Write(7, "7")

	mem.Write(8, "8")
	mem.Write(9, "9")
	mem.Write(10, "10")
	mem.Write(11, "11")
	mem.Write(12, "12")
	mem.Write(13, "13")
	mem.Write(14, "14")

	out, ok := mem.GetAll()
	if !ok {
		fmt.Println("no elements")
		return
	}

	for _, i := range out {
		fmt.Println(i.(string))
	}
}
