package main

import (
	"context"
	"flag"
	"fmt"

	//"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"rainier/internal/config"

	"rainier/internal/handlers"

	"rainier/internal/metrics"

	"rainier/internal/middleware"

	"rainier/internal/routes"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	fmt.Println("== Starting Service ==")

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Can not get configuration")
	}

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().Str("service", "warmachine").Msg("Starting up")
	log.Info().Str("service", "warmachine").Msgf("USING SITE: %s", cfg.Site)
	log.Info().Str("service", "warmachine").Msgf("ENV: %s", cfg.Env)
	log.Info().Str("service", "warmachine").Msgf("SEARCHLAYER ENDPOINT: %s", cfg.SearchLayer)
	log.Info().Str("service", "warmachine").Msgf("LOCATION HELPER SERVICE ENDPOINT: %s", cfg.LocationHelper)

	log.Info().Msgf("== Initializing Configuration ==")

	hctx := handlers.CTXHandlerContext(&cfg)
	mwctx := middleware.CTXMiddlewareContext(&cfg)

	middleware := metrics.NewPrometheusMiddleware(metrics.Opts{})

	r := routes.GetRoutes(hctx, mwctx)

	r.Handle("/metrics", promhttp.Handler())

	// This instruments all the routes, probably dont want this for a large application with lots of potential routes. Example: /name/john-oconnor
	r.Use(middleware.InstrumentHandlerDuration)

	log.Info().Str("service", "warmachine").Msg("Now happily serving requests")

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:5001",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("Error starting up HTTP Listener")
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.

	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info().Msg("Shutting down server")
	os.Exit(0)
}
