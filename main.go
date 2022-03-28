package main

import "fmt"

func main() {
	if err := run(); err != nil {
		panic(err)
	}

}

func run() error {
	fmt.Println("Hello!")
	return nil
}
