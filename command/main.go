package main

import "os"

func main() {
	if err := RunMain(os.Args); err != nil {
		os.Exit(1)
	}
}
