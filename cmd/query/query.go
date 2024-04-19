package query

import (
	"fmt"

	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/spf13/cobra"
)

func NewQueryCommand(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "query",
		Short: "Run the configured queries and output the results",
		RunE: func(cmd *cobra.Command, args []string) error {
			results := app.Collector.Collect(cmd.Context())

			for _, result := range results {
				fmt.Printf("Service '%s'\n", result.Service.Name)
				status := "down"
				if result.Status {
					status = "up"
				}
				if !result.Success {
					status = "unknown"
				}
				fmt.Printf("  Status: %s\n", status)
				fmt.Printf("  Scrape successful: %t\n", result.Success)
			}

			return nil
		},
	}
}
