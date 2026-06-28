package main

import (
	"fmt"
	"os"
	"vessel/internal/container"
)

func main() {
	var err error

    switch os.Args[1] {
    case "run":
        err = container.Run()
    case "child":
        err = container.Child()
	default :
		fmt.Println("invalid command")
		os.Exit(1)
    }

	if err != nil {
		fmt.Fprintf(os.Stderr, "init container failed: %v\n", err)
		os.Exit(1)
	}
}
