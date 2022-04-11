package main

import (
	"tupulung/config"
	"tupulung/utilities"
)

func main() {
	config := config.Get()
	utilities.NewMysqlGorm(config)
}