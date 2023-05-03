package main

import (
	"fmt"
)

func main() {
	releases := Releases()
	fmt.Println(len(releases))
}
