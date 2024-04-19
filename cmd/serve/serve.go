package serve

import (
	"context"
	"errors"
	"fmt"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/henrywhitaker3/prompage/internal/config"
	"github.com/henrywhitaker3/prompage/internal/http"
	"github.com/spf13/cobra"
)

func NewServeCommand(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run the status page http server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigs
				fmt.Println("Received interrupt, stopping...")
				cancel()
			}()

			http := http.NewHttp(conf)

			go func() {
				if err := http.Serve(); err != nil {
					if !errors.Is(err, stdhttp.ErrServerClosed) {
						fmt.Println(fmt.Errorf("http server failed: %v", err))
						cancel()
					}
				}
			}()

			<-ctx.Done()
			return http.Stop(context.Background())
		},
	}
}
