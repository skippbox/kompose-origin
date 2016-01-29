package main

import (
	"os"

	"github.com/codegangsta/cli"
	cliApp "github.com/docker/libcompose/cli/app"
	"github.com/docker/libcompose/cli/command"
	dockerApp "github.com/docker/libcompose/cli/docker/app"
	"github.com/docker/libcompose/version"
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
	app.Flags = append(command.CommonFlags(), dockerApp.DockerClientFlags()...)
	app.Commands = []cli.Command{
		command.BuildCommand(factory),
		command.CreateCommand(factory),
		command.UpCommand(factory),
		command.StartCommand(factory),
		command.LogsCommand(factory),
		command.RestartCommand(factory),
		command.StopCommand(factory),
		command.ScaleCommand(factory),
		command.RmCommand(factory),
		command.PullCommand(factory),
		command.KillCommand(factory),
		command.PortCommand(factory),
		command.KuberCommand(factory),
		command.PsCommand(factory),
		command.KuberConfigCommand(factory),
		command.PauseCommand(factory),
		command.UnpauseCommand(factory),
	}

	app.Run(os.Args)
}
