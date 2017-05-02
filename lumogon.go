package main

import (
	"fmt"
	"os"

	"github.com/puppetlabs/transparent-containers/cli/capabilities"
	"github.com/puppetlabs/transparent-containers/cli/cmd"
)

func main() {
	capabilities.Init()
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
