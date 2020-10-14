package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/spf13/cobra"

	"ticket-reservation/app"
	"ticket-reservation/http_api"
	"ticket-reservation/log"
	"ticket-reservation/prometheus"
)

func serveHTTPAPI(ctx context.Context, logger log.Logger, httpAPI *http_api.API) {
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	router := httpAPI.Init()

	s := &http.Server{
		Addr:        fmt.Sprintf(":%d", httpAPI.Config.Port),
		Handler:     cors(router),
		ReadTimeout: 2 * time.Minute,
	}

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			logger.Errorf("%+v", err)
		}
		close(done)
	}()

	logger.Infof("Serving HTTP API at http://127.0.0.1:%d", httpAPI.Config.Port)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logger.Errorf("%+v", err)
	}
	<-done
}

var serveAPICmd = &cobra.Command{
	Use:   "serve-api",
	Short: "Start API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, err := getLogger()
		if err != nil {
			return err
		}

		app, err := app.New(logger, nil)
		if err != nil {
			return err
		}

		httpAPI, err := http_api.New(app)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			var count int
			c := make(chan os.Signal, 2)
			signal.Notify(c, syscall.SIGTERM, os.Interrupt)
			go func() {
				for {
					select {
					case <-c:
						count++
						if count == 2 {
							logger.Infof("Forcefully exiting...")
							os.Exit(1)
						}
						logger.Infof("Signal SIGKILL caught. shutting down...")
						logger.Infof("Catching SIGKILL one more time will forcefully exit")
						cancel()
					}
				}
			}()
		}()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			serveHTTPAPI(ctx, logger, httpAPI)
			app.Close()
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer cancel()
			prometheus.ServePrometheusMetrics(ctx, logger)
		}()

		wg.Wait()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveAPICmd)
}
