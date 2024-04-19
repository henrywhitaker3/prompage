package cmd

import (
	"github.com/henrywhitaker3/prompage/cmd/query"
	"github.com/henrywhitaker3/prompage/cmd/serve"
	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prompage",
		Short: "PromPage Status Page",
	}
}

func LoadSubCommands(cmd *cobra.Command, app *app.App) {
	cmd.AddCommand(serve.NewServeCommand(app))
	cmd.AddCommand(query.NewQueryCommand(app))
}
