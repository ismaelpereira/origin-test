package main

import (
	"origin-challenge/api"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}

}

func run() error {

	api.ApiHandler()
	return nil
}
