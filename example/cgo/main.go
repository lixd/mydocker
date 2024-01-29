package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
)

func main() {
	Seed(123)
	// Output：Random:  128959393
	fmt.Println("Random: ", Random())
}

// Seed 初始化随机数产生器
func Seed(i int) {
	C.srandom(C.uint(i))
}

// Random 产生一个随机数
func Random() int {
	return int(C.random())
}
