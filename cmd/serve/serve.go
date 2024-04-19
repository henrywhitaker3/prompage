package serve

import (
	"fmt"

	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/spf13/cobra"
)

func NewServeCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run the status page http server",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(conf.Port)
			return nil
		},
	}
}
