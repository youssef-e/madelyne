package main

import (
	"fmt"
	"github.com/madelyne-io/madelyne/tester"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("You must provide a valid config file")
		os.Exit(1)
	}

	suite, err := tester.Load(os.Args[1])
	if err != nil {
		fmt.Println("Cannot read config file : ", err)
		os.Exit(2)
	}
	fmt.Println("Testing REST API with Madelyne")
	err = suite.Run()
	if err != nil {
		fmt.Println("\n\nError while running test: ", err)
		os.Exit(3)
	}
	fmt.Println("Success")
}
