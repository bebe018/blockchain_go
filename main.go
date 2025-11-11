package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("os.Args:", os.Args)
	cli := CLI{}
	cli.Run()
}
