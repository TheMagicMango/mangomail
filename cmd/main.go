package main

import (
	"os"

	"github.com/TheMagicMango/mangomail/cmd/root"
)

func main() {
	err := root.Cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
