package prometheus

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"ticket-reservation/log"
)

const prometheusNamespace = "ticket_reservation"

var config *Config

func serveHTTPMetrics(ctx context.Context, logger log.Logger, port int) {
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	router := mux.NewRouter()
	router.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	s := &http.Server{
		Addr:        fmt.Sprintf(":%d", port),
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

	logger.Infof("Serving HTTP metrics API at port: %d", port)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logger.Errorf("%+v", err)
	}
	<-done
}

func ServePrometheusMetrics(ctx context.Context, logger log.Logger) (err error) {
	internalLogger := logger.WithFields(log.Fields{
		"module": "prometheus",
	})

	config, err = InitConfig()
	if err != nil {
		return err
	}

	if !config.Enable {
		return nil
	}

	register()

	// Start HTTP API serving metrics
	serveHTTPMetrics(ctx, internalLogger, config.MetricsPort)

	return nil
}

func register() {
	prometheus.MustRegister(httpAPICallDurationHistogram)
}
