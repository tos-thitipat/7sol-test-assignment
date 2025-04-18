package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"pie_fire_dine/config"
	"pie_fire_dine/external"
	"pie_fire_dine/handler"
	"pie_fire_dine/repository"
	"pie_fire_dine/service"

	"github.com/rs/cors"
)

func Start() {
	meatRepo := repository.NewMeatRepository(config.GetCsvMeatCategoryPath())

	extRequest := external.NewHttpRequester(external.RequestHttp)
	meatSummaryService := service.NewMeatSummaryService(meatRepo, extRequest)

	pieFireDineHandler := handler.NewMeatSummaryHandler(meatSummaryService)

	c := cors.New(cors.Options{
		AllowedHeaders: []string{"*"},
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET"},
	})
	handler := c.Handler(Router(pieFireDineHandler))
	startServer(handler)
}

func startServer(handler http.Handler) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	slog.Info(fmt.Sprintf("starting %s server on port: %s", config.GetAppName(), config.GetPort()))

	server := &http.Server{
		Addr:    ":" + config.GetPort(),
		Handler: handler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	slog.Info(fmt.Sprintf("shutting %s server down gracefully, press Ctrl+C again to force", config.GetAppName()))

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}
