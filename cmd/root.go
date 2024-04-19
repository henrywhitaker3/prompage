package cmd

import (
	"github.com/henrywhitaker3/prompage/cmd/serve"
	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Short: "PromPage Status Page",
	}
}

func LoadSubCommands(cmd *cobra.Command, conf *config.Config) {
	cmd.AddCommand(serve.NewServeCommand(conf))

}
