package main

import (
	"context"

	"github.com/jainabhishek5986/employee-records/cmd"
	"github.com/jainabhishek5986/employee-records/pkg/errs"
	"github.com/jainabhishek5986/employee-records/pkg/global"
	"github.com/jainabhishek5986/employee-records/pkg/zaplogger"
	"go.uber.org/zap"
)

// Main function just executes root command `ts` and this is the entry point
// this project structure is inspired from `cobra` package
func main() {
	// Initiate Logger with Log File Name
	err := zaplogger.InitLogger(global.LogFileName)
	if err != nil {
		zaplogger.Fatal(context.Background(), errs.InitiateLoggerError, zap.Error(err))
	}
	if err := cmd.RootCommand().Execute(); err != nil {
		zaplogger.Fatal(context.Background(), errs.StartServerError, zap.Error(err))
	}
}
