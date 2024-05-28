package main

import (
	"fmt"
	"os"

	"github.com/henrywhitaker3/prompage/cmd"
	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/henrywhitaker3/prompage/internal/querier"
	"github.com/henrywhitaker3/prompage/internal/resources/views"
	"github.com/spf13/pflag"
)

var (
	configPath string
	version    string
)

//go:generate npm run build

func main() {
	pflag.StringVarP(&configPath, "config", "c", "prompage.yaml", "The location of the config file")

	root := cmd.NewRootCmd()

	pflag.Parse()

	conf, err := config.Load(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	qs, err := querier.BuildQueriers(conf.Datasources)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	app := app.NewApp(conf, qs)
	app.Version = version

	views.MustCompile()

	cmd.LoadSubCommands(root, app)

	if err := root.Execute(); err != nil {
		os.Exit(2)
	}
}
