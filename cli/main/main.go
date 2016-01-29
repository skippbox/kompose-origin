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
		command.EventsCommand(factory),
		command.DownCommand(factory),
		command.KillCommand(factory),
		command.LogsCommand(factory),
		command.PauseCommand(factory),
		command.PortCommand(factory),
		command.PsCommand(factory),
		command.PullCommand(factory),
		command.RestartCommand(factory),
		command.RmCommand(factory),
		command.RunCommand(factory),
		command.ScaleCommand(factory),
		command.StartCommand(factory),
		command.StopCommand(factory),
		command.UnpauseCommand(factory),
		command.UpCommand(factory),
		command.VersionCommand(factory),
		command.KuberCommand(factory),
		command.PsCommand(factory),
		command.KuberConfigCommand(factory),
	}

	app.Run(os.Args)

}
