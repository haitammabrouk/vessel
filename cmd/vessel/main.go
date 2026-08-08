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
		containerId := os.Args[2]
        err = container.Child(containerId)
	case "ps":
        err = container.ListContainers()
	case "stop":
		containerId := os.Args[2]
        err = container.StopContainer(containerId)
	default :
		fmt.Println("invalid command")
		os.Exit(1)
    }

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
