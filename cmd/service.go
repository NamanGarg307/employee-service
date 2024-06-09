package cmd

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/jainabhishek5986/employee-records/config"
	"github.com/jainabhishek5986/employee-records/pkg/global"
	"github.com/jainabhishek5986/employee-records/pkg/models"
	"github.com/jainabhishek5986/employee-records/pkg/transport/http"
	"github.com/jainabhishek5986/employee-records/pkg/waitgroup"
	"github.com/jainabhishek5986/employee-records/pkg/zaplogger"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	// use mysql connector
	_ "github.com/go-sql-driver/mysql" // use mysql connector
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// RootCommand will setup and return the root command
func RootCommand() *cobra.Command {
	rootCmd := cobra.Command{
		Use:     "employee-records-service",
		Long:    `Microservice for Employee Records`,
		Version: global.BinaryVersion,
		Run:     run,
	}
	// define flags used for this command
	err := AttachCLIFlags(&rootCmd)
	if err != nil {
		zaplogger.Fatal(context.Background(), `Something went wrong while getting attached 
				flags`, zap.Error(err))
	}

	return &rootCmd
}

/*
run : This function start the process

Parameters
---------
cmd: Cobra command object
*/
func run(cmd *cobra.Command, args []string) {

	// getting Config object that have all the global
	// configuration information
	cfg, err := config.Load(cmd)
	if err != nil {
		zaplogger.Fatal(context.Background(), "Failed to load config", zap.Error(err))
	}

	// create a context that we can cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// timeout in seconds
	const GracefulTimeout = 40000 * time.Millisecond

	// load environment variables from .env if available
	err = godotenv.Load()
	if err != nil {
		zaplogger.Warn(ctx, "Warning: Environment file not found")
	}

	db := DBConnection(ctx)

	waitgroup.Gwg.Add(1)
	go func() {
		defer waitgroup.Gwg.Done()
		// setup http server
		err = http.Setup(ctx, cfg, &waitgroup.Gwg, db)
		if err != nil {
			zaplogger.Error(ctx, "Something Went Wrong", zap.Error(err))
		}
	}()

	waitgroup.Gwg.Add(1)

	// listen for C-c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// create channel to mark status of waitgroup
	// this is required to brutally kill application in case of
	// timeout
	done := make(chan struct{})

	// asynchronously wait for all the go routines
	go func() {
		// and wait for all go routines
		waitgroup.Gwg.Wait()
		zaplogger.Debug(ctx, "main: all goroutines have finished.")
		close(done)
	}()

	// wait for signal channel
	select {
	case <-c:

		zaplogger.Debug(ctx, "main: received C-c - attempting graceful shutdown ....")
		// tell the goroutines to stop
		zaplogger.Debug(ctx, "main: telling goroutines to stop")
		cancel()
		select {
		case <-done:
			zaplogger.Debug(ctx, "Go routines exited within timeout")
		case <-time.After(GracefulTimeout):
			zaplogger.Error(ctx, "Graceful timeout exceeded. Brutally killing the application")
		}

	case <-done:
		defer os.Exit(0)
	}
}

/*
DBConnection: Making db connection using gorm

Parameters
-----------
cfg: Config Object
logger: Logging object

Return
---------
db: DB connection object
*/
func DBConnection(ctx context.Context) (db *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		zaplogger.Panic(ctx, "Unable to connect to db. Exiting", zap.Error(err))
	}
	err = db.AutoMigrate(models.Employee{})
	if err != nil {
		zaplogger.Panic(ctx, "Unable to run migrations to db. Exiting", zap.Error(err))
	}
	return db

}
