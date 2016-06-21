package main

import (
	"os"

	"github.com/codegangsta/cli"
	cliApp "github.com/skippbox/kompose/cli/app"
	"github.com/skippbox/kompose/cli/command"
	dockerApp "github.com/skippbox/kompose/cli/docker/app"
	"github.com/skippbox/kompose/version"
)

func main() {
	factory := &dockerApp.ProjectFactory{}

	app := cli.NewApp()
	app.Name = "kompose"
	app.Usage = "Command line interface for Skippbox."
	app.Version = version.VERSION + " (" + version.GITCOMMIT + ")"
	app.Author = "Skippbox Compose Contributors"
	app.Email = "https://github.com/skippbox/kompose"
	app.Before = cliApp.BeforeApp
	app.Flags = append(command.CommonFlags())
	app.Commands = []cli.Command{
		command.ConvertCommand(factory),
		command.UpCommand(factory),
		command.PsCommand(factory),
		command.DeleteCommand(factory),
		command.ScaleCommand(factory),
	}

	app.Run(os.Args)
}

