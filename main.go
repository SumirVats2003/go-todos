package main

import "github.com/SumirVats2003/go-todo/cmd"

func main() {
	db := cmd.InitApp()
	defer db.Close()
}
