//go:generate statik -src=./public

package main

import (
	"github.com/uphy/procman/cli"
	_ "github.com/uphy/procman/statik"
)

func main() {
	cli.Execute()
}
