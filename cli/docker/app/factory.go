package app

import (
	"github.com/codegangsta/cli"
	"github.com/skippbox/kompose/cli/command"
	"github.com/skippbox/kompose/cli/logger"
	"github.com/skippbox/kompose/docker"
	"github.com/skippbox/kompose/project"
)

// ProjectFactory is a struct that hold the app.ProjectFactory implementation.
type ProjectFactory struct {
}

// Create implements ProjectFactory.Create using docker client.
func (p *ProjectFactory) Create(c *cli.Context) (*project.Project, error) {
	context := &docker.Context{}
	context.LoggerFactory = logger.NewColorLoggerFactory()
	Populate(context, c)
	command.Populate(&context.Context, c)

	return docker.NewProject(context)
}

