package main

import (
	"fmt"
	"os"

	"github.com/puppetlabs/lumogon/capabilities"
	"github.com/puppetlabs/lumogon/cmd"
)

func main() {
	capabilities.Init()
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
