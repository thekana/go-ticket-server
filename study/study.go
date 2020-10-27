package main

import (
	"fmt"
	"github.com/pkg/errors"
)

func main() {
	fmt.Println(testdefer().Error())
}

func testdefer() error {
	err := errors.New("Start")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recover")
		} else if err != nil {
			fmt.Println("else if")
			err = errors.Wrap(err, "Got err")
		} else {
			fmt.Println("else")
			err = errors.New("No err")
		}
	}()
	err = nil
	return errors.New("Static")
}
