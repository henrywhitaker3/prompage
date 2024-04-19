package main

import (
	"fmt"
	"os"

	"github.com/henrywhitaker3/prompage/cmd"
	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/spf13/pflag"
)

var (
	configPath string
)

func main() {
	pflag.StringVarP(&configPath, "config", "c", "prompage.yaml", "The location of the config file")

	root := cmd.NewRootCmd()

	pflag.Parse()

	conf, err := config.Load(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cmd.LoadSubCommands(root, conf)

	if err := root.Execute(); err != nil {
		os.Exit(2)
	}
}
