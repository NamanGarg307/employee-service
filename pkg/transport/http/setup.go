package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jainabhishek5986/employee-records/pkg/errs"
	"github.com/jainabhishek5986/employee-records/pkg/global"
	"github.com/jainabhishek5986/employee-records/pkg/zaplogger"
	"github.com/lestrrat-go/backoff"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/jainabhishek5986/employee-records/config"
)

/*
StartAPIServer : start the API server and call of functions like
1- middleware for global logging
2- Setting up cors
3- Routes function to get all the endpoint
4- Register API routes

Parameters
----------
ctx: Global context
config: Config object
wg: Wait group object
db: Database connection
logger: Global logging object
*/
func StartAPIServer(ctx context.Context, conf *config.Config,
	wg *sync.WaitGroup, db *gorm.DB) error {

	wg.Add(1)
	defer wg.Done()

	// Set gin to release mode
	gin.SetMode(gin.ReleaseMode)

	zaplogger.Info(ctx, "Setting up http handler")
	router := gin.Default()

	// Recovery middleware recovers from any panics and writes a 500
	// if there was one
	router.Use(gin.Recovery())

	errChan := make(chan error)

	// All the router groups
	v1RoutesGroup := router.Group("/api/v1")

	// Cors config for rest of the routes
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	v1RoutesGroup.Use(cors.New(corsConfig))

	// Registering API Routes
	RegisterAPIRoutes(v1RoutesGroup, db)

	// HTTP server instance
	srv := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: router,
	}

	// channel to signal server process exit
	done := make(chan struct{})

	go func() {
		zaplogger.Info(ctx, "Starting server on port", zap.String("port", conf.Port))
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zaplogger.Error(ctx, "listen", zap.Error(err))
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		const GracefulTimeout = 20000 * time.Millisecond
		shutdownCtx, cancel := context.WithTimeout(context.Background(), GracefulTimeout)

		defer cancel()
		zaplogger.Info(shutdownCtx, "Caller has requested graceful shutdown. shutting down the server")
		if err := srv.Shutdown(shutdownCtx); err != nil {
			zaplogger.Error(shutdownCtx, fmt.Sprintf("Server Shutdown: error -%s", err.Error()))
		}
		return nil
	case err := <-errChan:
		return err
	case <-done:
		return nil
	}
}

/*
Setup function just set up the process to start the API server

Parameters
----------
config: Config object
ctx: Global context
wg: Wait group object
db: DB object
*/
func Setup(ctx context.Context, conf *config.Config, wg *sync.WaitGroup, db *gorm.DB) error {
	zaplogger.Info(ctx, "Starting API server")

	var policy = backoff.NewExponential(
		backoff.WithInterval(global.FiveHundred*time.Millisecond), // base interval
		backoff.WithJitterFactor(global.PointZeroFive),            // 5% jitter
		backoff.WithMaxRetries(global.MaxAPIServerStartAttempts),  // If not specified, default number of retries is 10
	)

	b, cancel := policy.Start(context.Background())
	defer cancel()

	for backoff.Continue(b) {
		// check for context
		select {
		case <-ctx.Done():
			zaplogger.Debug(ctx, "Context cancelled. Stopping sink proxy")
			return nil
		default:
			err := StartAPIServer(ctx, conf, wg, db)
			if err != nil {
				zaplogger.Error(ctx, errs.APIServerStartError, zap.Error(err))
			} else {
				return nil
			}
		}
	}

	// in the case of retry attempts to exceed, return with error
	return errors.New("failed to start server after maximum retries")
}
