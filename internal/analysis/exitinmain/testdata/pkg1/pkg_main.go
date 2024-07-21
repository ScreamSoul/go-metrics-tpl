package main

import (
	"fmt"
	"os"
)

func main() {

	fmt.Printf("Some log")
	os.Exit(1) // want `use of os.Exit\(\) in main function is discouraged`
}
