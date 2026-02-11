// (c) Magic Mango and individual authors
// SPDX-License-Identifier: Apache-2.0

// This package contains the mangomail CLI binary.
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
