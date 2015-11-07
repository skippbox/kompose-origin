package app

import (
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
	"io/ioutil"
	"path/filepath"

	"golang.org/x/net/context"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/pkg/api/unversioned"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

// ProjectAction is an adapter to allow the use of ordinary functions as libcompose actions.
// Any function that has the appropriate signature can be register as an action on a codegansta/cli command.
//
// cli.Command{
//		Name:   "ps",
//		Usage:  "List containers",
//		Action: app.WithProject(factory, app.ProjectPs),
//	}
type ProjectAction func(project project.APIProject, c *cli.Context) error

// BeforeApp is an action that is executed before any cli command.
func BeforeApp(c *cli.Context) error {
	if c.GlobalBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.Warning("Note: This is an experimental alternate implementation of the Compose CLI (https://github.com/docker/compose)")
	return nil
}

// WithProject is a helper function to create a cli.Command action with a ProjectFactory.
func WithProject(factory ProjectFactory, action ProjectAction) func(context *cli.Context) error {
	return func(context *cli.Context) error {
		p, err := factory.Create(context)
		if err != nil {
			logrus.Fatalf("Failed to read project: %v", err)
		}
		return action(p, context)
	}
}

// ProjectPs lists the containers.
func ProjectPs(p project.APIProject, c *cli.Context) error {
	qFlag := c.Bool("q")
	allInfo, err := p.Ps(context.Background(), qFlag, c.Args()...)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	os.Stdout.WriteString(allInfo.String(!qFlag))
	return nil
}

func ProjectKuberConfig(p *project.Project, c *cli.Context) {
	url := c.String("host")

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	confDir := usr.HomeDir
    if err != nil {
            log.Fatal(err)
    }
	outputFileName := fmt.Sprintf(".kuberconfig")
	outputFilePath := filepath.Join(confDir, outputFileName)
	if _, err :=  os.Stat(outputFilePath); os.IsNotExist(err) {
		f, err := os.Create(outputFilePath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}
	wurl := []byte(url)
	if err := ioutil.WriteFile(outputFilePath, wurl, 0644); err != nil {
		logrus.Fatalf("Failed to write k8s api server address to %s: %v", outputFilePath, err)
	}
}

func ProjectKuber(p *project.Project, c *cli.Context) {
	//outputDir := c.String("output")
	composeFile := c.String("file")

	p = project.NewProject(&project.Context{
		ProjectName: "kube",
		ComposeFile: composeFile,
	})

	if err := p.Parse(); err != nil {
		logrus.Fatalf("Failed to parse the compose project from %s: %v", composeFile, err)
	}
	//if err := os.MkdirAll(outputDir, 0755); err != nil {
	//	logrus.Fatalf("Failed to create the output directory %s: %v", outputDir, err)
	//}

	//Get config client
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	confDir := usr.HomeDir
    if err != nil {
            log.Fatal(err)
    }
	outputFileName := fmt.Sprintf(".kuberconfig")
	outputFilePath := filepath.Join(confDir, outputFileName)
	readServer, readErr := ioutil.ReadFile(outputFilePath)
	var server string = "127.0.0.1"

	if readErr == nil {
		// logrus.Fatalf("Failed to read k8s api server address from %s: %v", outputFilePath, readErr)
		server = string(readServer)
	}
	// server := string(readServer)
	// if server == "" {
	// 	logrus.Fatalf("K8s api server address isn't defined in %s", outputFilePath)
	// }

	var mServices map[string]api.Service = make(map[string]api.Service)
	var serviceLinks []string

	version := "v1"
	// create new client
	client := client.NewOrDie(&client.Config{Host: server, Version: version})
	for name, service := range p.Configs {
		rc := &api.ReplicationController{
			TypeMeta: unversioned.TypeMeta{
				Kind:       "ReplicationController",
				APIVersion: "v1",
			},
			ObjectMeta: api.ObjectMeta{
				Name:   name,
				Labels: map[string]string{"service": name},
			},
			Spec: api.ReplicationControllerSpec{
				Replicas: 1,
				Selector: map[string]string{"service": name},
				Template: &api.PodTemplateSpec{
					ObjectMeta: api.ObjectMeta{
						Labels: map[string]string{"service": name},
					},
					Spec: api.PodSpec{
						Containers: []api.Container{
							{
								Name:  name,
								Image: service.Image,
							},
						},
					},
				},
			},
		}
		sc := &api.Service{
			TypeMeta: unversioned.TypeMeta{
				Kind:       "Service",
				APIVersion: "v1",
			},
			ObjectMeta: api.ObjectMeta{
				Name:   name,
				Labels: map[string]string{"service": name},
			},
			Spec: api.ServiceSpec{
				Selector: map[string]string{"service": name},
			},
		}
		
		// Configure the container ports.
		var ports []api.ContainerPort
		for _, port := range service.Ports {
			var character string = ":"
			if strings.Contains(port, character) {
				//portNumber := port[0:strings.Index(port, character)]
				targetPortNumber := port[strings.Index(port, character) + 1: len(port)]
				targetPortNumberInt, err := strconv.Atoi(targetPortNumber)
				if err != nil {
					logrus.Fatalf("Invalid container port %s for service %s", port, name)
				}
				ports = append(ports, api.ContainerPort{ContainerPort: targetPortNumberInt})
			} else {
				portNumber, err := strconv.Atoi(port)
				if err != nil {
					logrus.Fatalf("Invalid container port %s for service %s", port, name)
				}
				ports = append(ports, api.ContainerPort{ContainerPort: portNumber})
			}			
		}

		rc.Spec.Template.Spec.Containers[0].Ports = ports

		// Configure the service ports.
		var servicePorts []api.ServicePort
		for _, port := range service.Ports {
			var character string = ":"
			if strings.Contains(port, character) {
				portNumber := port[0:strings.Index(port, character)]
				targetPortNumber := port[strings.Index(port, character) + 1: len(port)]
				portNumberInt, err := strconv.Atoi(portNumber)				
				if err != nil {
					logrus.Fatalf("Invalid container port %s for service %s", port, name)
				}
				targetPortNumberInt, err1 := strconv.Atoi(targetPortNumber)
				if err1 != nil {
					logrus.Fatalf("Invalid container port %s for service %s", port, name)
				}
				var targetPort util.IntOrString
				targetPort.StrVal = targetPortNumber
				targetPort.IntVal = targetPortNumberInt
				servicePorts = append(servicePorts, api.ServicePort{Port: portNumberInt, Name: portNumber, Protocol: "TCP", TargetPort: targetPort})	
			} else {
				portNumber, err := strconv.Atoi(port)
				if err != nil {
					logrus.Fatalf("Invalid container port %s for service %s", port, name)
				}
				var targetPort util.IntOrString
				targetPort.StrVal = strconv.Itoa(portNumber)
				targetPort.IntVal = portNumber
				servicePorts = append(servicePorts, api.ServicePort{Port: portNumber, Name: strconv.Itoa(portNumber), Protocol: "TCP", TargetPort: targetPort})
			}			
		}
		sc.Spec.Ports = servicePorts

		// Configure the container restart policy.
		switch service.Restart {
		case "", "always":
			rc.Spec.Template.Spec.RestartPolicy = api.RestartPolicyAlways
		case "no":
			rc.Spec.Template.Spec.RestartPolicy = api.RestartPolicyNever
		case "on-failure":
			rc.Spec.Template.Spec.RestartPolicy = api.RestartPolicyOnFailure
		default:
			logrus.Fatalf("Unknown restart policy %s for service %s", service.Restart, name)
		}

		datarc, err := json.MarshalIndent(rc, "", "  ")
		if err != nil {
			logrus.Fatalf("Failed to marshal the replication controller: %v", err)
		}
		fmt.Println(string(datarc))

		datasc, er := json.MarshalIndent(sc, "", "	")
		if er != nil {
			logrus.Fatalf("Failed to marshal the service controller: %v", er)
		}
		fmt.Println(string(datasc))

		mServices[name] = *sc			

		if len(service.Links.Slice()) > 0 {
			for i := 0; i < len(service.Links.Slice()); i++ {
				var data string = service.Links.Slice()[i]
				if len(serviceLinks) == 0 {
					serviceLinks = append(serviceLinks, data)
				} else {
					for _, v := range serviceLinks {
						if v != data {
							serviceLinks = append(serviceLinks, data)
						}						
					}
				}
			}			
			
		}
				
		// call create RC api
		rcCreated, err := client.ReplicationControllers(api.NamespaceDefault).Create(rc)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(rcCreated)

		for k, v := range mServices {
			fmt.Println(k)
			for i :=0; i < len(serviceLinks); i++ {
				if serviceLinks[i] == k {
					// call create SVC api
					scCreated, err := client.Services(api.NamespaceDefault).Create(&v)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(scCreated)
				}
			}
		}

		//Test


	}
}

func ProjectPort(p *project.Project, c *cli.Context) {
	if len(c.Args()) != 2 {
		return cli.NewExitError("Please pass arguments in the form: SERVICE PORT", 1)
	}

	index := c.Int("index")
	protocol := c.String("protocol")
	serviceName := c.Args()[0]
	privatePort := c.Args()[1]

	port, err := p.Port(context.Background(), index, protocol, serviceName, privatePort)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	fmt.Println(port)
	return nil
}

// ProjectStop stops all services.
func ProjectStop(p project.APIProject, c *cli.Context) error {
	err := p.Stop(context.Background(), c.Int("timeout"), c.Args()...)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// ProjectDown brings all services down (stops and clean containers).
func ProjectDown(p project.APIProject, c *cli.Context) error {
	options := options.Down{
		RemoveVolume:  c.Bool("volumes"),
		RemoveImages:  options.ImageType(c.String("rmi")),
		RemoveOrphans: c.Bool("remove-orphans"),
	}
	err := p.Down(context.Background(), options, c.Args()...)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// ProjectBuild builds or rebuilds services.
func ProjectBuild(p project.APIProject, c *cli.Context) error {
	config := options.Build{
		NoCache:     c.Bool("no-cache"),
		ForceRemove: c.Bool("force-rm"),
		Pull:        c.Bool("pull"),
	}
	err := p.Build(context.Background(), config, c.Args()...)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// ProjectCreate creates all services but do not start them.
func ProjectCreate(p project.APIProject, c *cli.Context) error {
	options := options.Create{
		NoRecreate:    c.Bool("no-recreate"),
		ForceRecreate: c.Bool("force-recreate"),
		NoBuild:       c.Bool("no-build"),
	}
	err := p.Create(context.Background(), options, c.Args()...)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// ProjectUp brings all services up.
func ProjectUp(p project.APIProject, c *cli.Context) error {
	options := options.Up{
		Create: options.Create{
			NoRecreate:    c.Bool("no-recreate"),
			ForceRecreate: c.Bool("force-recreate"),
			NoBuild:       c.Bool("no-build"),
		},
	}
	ctx, cancelFun := context.WithCancel(context.Background())
	err := p.Up(ctx, options, c.Args()...)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	if !c.Bool("d") {
		signalChan := make(chan os.Signal, 1)
		cleanupDone := make(chan bool)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		errChan := make(chan error)
		go func() {
			errChan <- p.Log(ctx, true, c.Args()...)
		}()
		go func() {
			select {
			case <-signalChan:
				fmt.Printf("\nGracefully stopping...\n")
				cancelFun()
				ProjectStop(p, c)
				cleanupDone <- true
			case err := <-errChan:
				if err != nil {
					logrus.Fatal(err)
				}
				cleanupDone <- true
			}
		}()
		<-cleanupDone
		return nil
	}
	return nil
}

// ProjectRun runs a given command within a service's container.
func ProjectRun(p project.APIProject, c *cli.Context) error {
	if len(c.Args()) == 1 {
		logrus.Fatal("No service specified")
	}

	serviceName := c.Args()[0]
	commandParts := c.Args()[1:]

	exitCode, err := p.Run(context.Background(), serviceName, commandParts)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return cli.NewExitError("", exitCode)
}

// ProjectStart starts services.
func ProjectStart(p project.APIProject, c *cli.Context) error {
	err := p.Start(context.Background(), c.Args()...)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// ProjectRestart restarts services.
func ProjectRestart(p project.APIProject, c *cli.Context) error {
	err := p.Restart(context.Background(), c.Int("timeout"), c.Args()...)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// ProjectLog gets services logs.
func ProjectLog(p project.APIProject, c *cli.Context) error {
	err := p.Log(context.Background(), c.Bool("follow"), c.Args()...)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// ProjectPull pulls images for services.
func ProjectPull(p project.APIProject, c *cli.Context) error {
	err := p.Pull(context.Background(), c.Args()...)
	if err != nil && !c.Bool("ignore-pull-failures") {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// ProjectDelete deletes services.
func ProjectDelete(p project.APIProject, c *cli.Context) error {
	options := options.Delete{
		RemoveVolume: c.Bool("v"),
	}
	if !c.Bool("force") {
		options.BeforeDeleteCallback = func(stoppedContainers []string) bool {
			fmt.Printf("Going to remove %v\nAre you sure? [yN]\n", strings.Join(stoppedContainers, ", "))
			var answer string
			_, err := fmt.Scanln(&answer)
			if err != nil {
				logrus.Error(err)
				return false
			}
			if answer != "y" && answer != "Y" {
				return false
			}
			return true
		}
	}
	err := p.Delete(context.Background(), options, c.Args()...)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// ProjectKill forces stop service containers.
func ProjectKill(p project.APIProject, c *cli.Context) error {
	err := p.Kill(context.Background(), c.String("signal"), c.Args()...)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// ProjectPause pauses service containers.
func ProjectPause(p project.APIProject, c *cli.Context) error {
	err := p.Pause(context.Background(), c.Args()...)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// ProjectUnpause unpauses service containers.
func ProjectUnpause(p project.APIProject, c *cli.Context) error {
	err := p.Unpause(context.Background(), c.Args()...)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

// ProjectScale scales services.
func ProjectScale(p project.APIProject, c *cli.Context) error {
	servicesScale := map[string]int{}
	for _, arg := range c.Args() {
		kv := strings.SplitN(arg, "=", 2)
		if len(kv) != 2 {
			return cli.NewExitError(fmt.Sprintf("Invalid scale parameter: %s", arg), 2)
		}

		name := kv[0]

		count, err := strconv.Atoi(kv[1])
		if err != nil {
			return cli.NewExitError(fmt.Sprintf("Invalid scale parameter: %v", err), 2)
		}

		servicesScale[name] = count
	}

	err := p.Scale(context.Background(), c.Int("timeout"), servicesScale)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}
