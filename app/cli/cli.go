package main

import (
	"go-as/app"
	"go-as/app/cli/commands"
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	var clis []*cli.Command

	container := app.BuildDIContainer()
	if err := container.Invoke(func(logger *zap.Logger) {
		handleError(container.Invoke(func(permissionsCli *commands.BoostrapPermissionsCLI) {
			clis = append(clis, &cli.Command{
				Name:   "BootstrapPermissions",
				Usage:  "Bootstrap the application permissions",
				Action: permissionsCli.Execute,
			})
		}), logger)
	}); err != nil {
		panic("Error trying to build command clis")
	}

	app := &cli.App{
		Commands: clis,
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func handleError(err error, logger *zap.Logger) {
	if err != nil {
		logger.Fatal(err.Error())
	}
}
