package main

import (
	"flag"
	"fmt"
)

func main() {
	command := flag.String("command", "", "init|generate|setup|status|up|down")
	arg := flag.String("arg", "", "additional info for command")
	flag.Parse()

	fmt.Println(command)
	fmt.Println(arg)
}
