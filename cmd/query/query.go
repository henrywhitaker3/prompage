package query

import (
	"fmt"

	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/henrywhitaker3/prompage/internal/querier"
	"github.com/spf13/cobra"
)

func NewQueryCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "query",
		Short: "Run the configured queries and output the results",
		RunE: func(cmd *cobra.Command, args []string) error {
			q, err := querier.NewQuerier(conf)
			if err != nil {
				return err
			}

			for _, service := range conf.Services {
				fmt.Printf("Querying for service '%s'\n", service.Name)
				for _, query := range service.Queries {
					res, err := q.Status(cmd.Context(), query)
					if err != nil {
						fmt.Println(fmt.Errorf("query service error: %v", err))
					}
					fmt.Printf("  - Query: '%s': %v\n", query.Name, res)
				}
			}

			return nil
		},
	}
}
