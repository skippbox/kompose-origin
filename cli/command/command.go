package command

import (
	"github.com/codegangsta/cli"
	"github.com/skippbox/kompose/cli/app"
	k8sApp "github.com/skippbox/kompose/cli/k8s/app"
	"github.com/skippbox/kompose/project"
)

// ConvertCommand defines the kompose convert subcommand.
func ConvertCommand(factory app.ProjectFactory) cli.Command {
	return cli.Command{
		Name:   "convert",
		Usage:  "Convert docker-compose.yml to Kubernetes objects",
		Action: app.WithProject(factory, k8sApp.ProjectKuberConvert),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "file,f",
				Usage:  "Specify an alternate compose file (default: docker-compose.yml)",
				Value:  "docker-compose.yml",
				EnvVar: "COMPOSE_FILE",
			},
			cli.BoolFlag{
				Name:  "deployment,d",
				Usage: "Generate a deployment resource file",
			},
			cli.BoolFlag{
				Name:  "daemonset,ds",
				Usage: "Generate a daemonset resource file",
			},
			cli.BoolFlag{
				Name:  "replicaset,rs",
				Usage: "Generate a replicaset resource file",
			},
			cli.BoolFlag{
				Name:  "chart,c",
				Usage: "Create a chart deployment",
			},
			cli.BoolFlag{
				Name:  "yaml, y",
				Usage: "Generate resource file in yaml format",
			},
		},
	}
}

// UpCommand defines the kompose up subcommand.
func UpCommand(factory app.ProjectFactory) cli.Command {
	return cli.Command{
		Name:   "up",
		Usage:  "Submit rc, svc objects to kubernetes API endpoint",
		Action: app.WithProject(factory, k8sApp.ProjectKuberUp),
	}
}

// PsCommand defines the kompose ps subcommand.
func PsCommand(factory app.ProjectFactory) cli.Command {
	return cli.Command{
		Name:   "ps",
		Usage:  "Get active data in the kubernetes cluster",
		Action: app.WithProject(factory, k8sApp.ProjectKuberPS),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "service,svc",
				Usage: "Get active services",
			},
			cli.BoolFlag{
				Name:  "replicationcontroller,rc",
				Usage: "Get active replication controller",
			},
		},
	}
}

// DeleteCommand defines the kompose delete subcommand.
func DeleteCommand(factory app.ProjectFactory) cli.Command {
	return cli.Command{
		Name:   "delete",
		Usage:  "Remove instantiated services/rc from kubernetes",
		Action: app.WithProject(factory, k8sApp.ProjectKuberDelete),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "replicationcontroller,rc",
				Usage: "Remove active replication controllers",
			},
			cli.BoolFlag{
				Name:  "service,svc",
				Usage: "Remove active services",
			},
			cli.StringFlag{
				Name:  "name",
				Usage: "Name of the object to remove",
			},
		},
	}
}

// ScaleCommand defines the kompose up subcommand.
func ScaleCommand(factory app.ProjectFactory) cli.Command {
	return cli.Command{
		Name:   "scale",
		Usage:  "Globally scale instantiated replication controllers",
		Action: app.WithProject(factory, k8sApp.ProjectKuberScale),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "scale",
				Usage: "New number of replicas",
			},
			cli.StringFlag{
				Name:  "replicationcontroller,rc",
				Usage: "A specific replication controller to scale",
			},
		},
	}
}

// CommonFlags defines the flags that are in common for all subcommands.
func CommonFlags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name: "verbose,debug",
		},
		cli.StringFlag{
			Name:   "file,f",
			Usage:  "Specify an alternate compose file (default: docker-compose.yml)",
			Value:  "docker-compose.yml",
			EnvVar: "COMPOSE_FILE",
		},
		cli.StringFlag{
			Name:  "project-name,p",
			Usage: "Specify an alternate project name (default: directory name)",
		},
	}
}

// Populate updates the specified project context based on command line arguments and subcommands.
func Populate(context *project.Context, c *cli.Context) {
	context.ComposeFile = c.GlobalString("file")
	context.ProjectName = c.GlobalString("project-name")

	if c.Command.Name == "logs" {
		context.Log = true
	} else if c.Command.Name == "up" {
		context.Log = !c.Bool("d")
		context.NoRecreate = c.Bool("no-recreate")
		context.ForceRecreate = c.Bool("force-recreate")
	} else if c.Command.Name == "scale" {
		context.Timeout = uint(c.Int("timeout"))
	}
}

