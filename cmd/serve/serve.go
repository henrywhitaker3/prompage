package serve

import (
	"context"
	"errors"
	"fmt"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/henrywhitaker3/prompage/internal/app"
	"github.com/henrywhitaker3/prompage/internal/http"
	"github.com/spf13/cobra"
)

func NewServeCommand(app *app.App) *cobra.Command {
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

			cache := http.NewResultCache(app)
			http := http.NewHttp(app)

			go cache.Work(ctx)
			go func() {
				if err := http.Start(); err != nil {
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
