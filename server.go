package main

import (
	"fmt"
	"tupulung/config"
)

func main() {
	config := config.Get()

	fmt.Println(config)
}