package app

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/libcompose/project"

	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/pkg/api/unversioned"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

type ProjectAction func(project *project.Project, c *cli.Context)

func BeforeApp(c *cli.Context) error {
	if c.GlobalBool("verbose") {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.Warning("Note: This is an experimental alternate implementation of the Compose CLI (https://github.com/docker/compose)")
	return nil
}

func WithProject(factory ProjectFactory, action ProjectAction) func(context *cli.Context) {
	return func(context *cli.Context) {
		p, err := factory.Create(context)
		if err != nil {
			log.Fatalf("Failed to read project: %v", err)
		}
		action(p, context)
	}
}

func ProjectPs(p *project.Project, c *cli.Context) {
	allInfo := project.InfoSet{}
	for name := range p.Configs {
		service, err := p.CreateService(name)
		if err != nil {
			logrus.Fatal(err)
		}

		info, err := service.Info()
		if err != nil {
			logrus.Fatal(err)
		}

		allInfo = append(allInfo, info...)
	}

	os.Stdout.WriteString(allInfo.String())
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
	if readErr != nil {
		logrus.Fatalf("Failed to read k8s api server address from %s: %v", outputFilePath, readErr)
	}
	server := string(readServer)
	if server == "" {
		logrus.Fatalf("K8s api server address isn't defined in %s", outputFilePath)
	}

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
			portNumber, err := strconv.Atoi(port)
			if err != nil {
				logrus.Fatalf("Invalid container port %s for service %s", port, name)
			}
			ports = append(ports, api.ContainerPort{ContainerPort: portNumber})
		}

		rc.Spec.Template.Spec.Containers[0].Ports = ports

		// Configure the service ports.
		var servicePorts []api.ServicePort
		for _, port := range service.Ports {
			portNumber, err := strconv.Atoi(port)
			if err != nil {
				logrus.Fatalf("Invalid container port %s for service %s", port, name)
			}
			var targetPort util.IntOrString
			targetPort.StrVal = strconv.Itoa(portNumber)
			targetPort.IntVal = portNumber
			servicePorts = append(servicePorts, api.ServicePort{Port: portNumber, Name: strconv.Itoa(portNumber), Protocol: "TCP", TargetPort: targetPort})
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
	}
}

func ProjectPort(p *project.Project, c *cli.Context) {
	if len(c.Args()) != 2 {
		logrus.Fatalf("Please pass arguments in the form: SERVICE PORT")
	}

	index := c.Int("index")
	protocol := c.String("protocol")

	service, err := p.CreateService(c.Args()[0])
	if err != nil {
		logrus.Fatal(err)
	}

	containers, err := service.Containers()
	if err != nil {
		logrus.Fatal(err)
	}

	if index < 1 || index > len(containers) {
		logrus.Fatalf("Invalid index %d", index)
	}

	output, err := containers[index-1].Port(fmt.Sprintf("%s/%s", c.Args()[1], protocol))
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Println(output)
}

func ProjectDown(p *project.Project, c *cli.Context) {
	err := p.Down(c.Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

func ProjectBuild(p *project.Project, c *cli.Context) {
	err := p.Build(c.Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

func ProjectCreate(p *project.Project, c *cli.Context) {
	err := p.Create(c.Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

func ProjectUp(p *project.Project, c *cli.Context) {
	err := p.Up(c.Args()...)
	if err != nil {
		logrus.Fatal(err)
	}

	if !c.Bool("d") {
		wait()
	}
}

func ProjectStart(p *project.Project, c *cli.Context) {
	err := p.Start(c.Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

func ProjectRestart(p *project.Project, c *cli.Context) {
	err := p.Restart(c.Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

func ProjectLog(p *project.Project, c *cli.Context) {
	err := p.Log(c.Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
	wait()
}

func ProjectPull(p *project.Project, c *cli.Context) {
	err := p.Pull(c.Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

func ProjectDelete(p *project.Project, c *cli.Context) {
	if !c.Bool("force") && len(c.Args()) == 0 {
		logrus.Fatal("Will not remove all services with out --force")
	}
	err := p.Delete(c.Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

func ProjectKill(p *project.Project, c *cli.Context) {
	err := p.Kill(c.Args()...)
	if err != nil {
		logrus.Fatal(err)
	}
}

func ProjectScale(p *project.Project, c *cli.Context) {
	// This code is a bit verbose but I wanted to parse everything up front
	order := make([]string, 0, 0)
	serviceScale := make(map[string]int)
	services := make(map[string]project.Service)

	for _, arg := range c.Args() {
		kv := strings.SplitN(arg, "=", 2)
		if len(kv) != 2 {
			logrus.Fatalf("Invalid scale parameter: %s", arg)
		}

		name := kv[0]

		count, err := strconv.Atoi(kv[1])
		if err != nil {
			logrus.Fatalf("Invalid scale parameter: %v", err)
		}

		if _, ok := p.Configs[name]; !ok {
			logrus.Fatalf("% is not defined in the template", name)
		}

		service, err := p.CreateService(name)
		if err != nil {
			logrus.Fatalf("Failed to lookup service: %s: %v", service, err)
		}

		order = append(order, name)
		serviceScale[name] = count
		services[name] = service
	}

	for _, name := range order {
		scale := serviceScale[name]
		logrus.Infof("Setting scale %s=%d...", name, scale)
		err := services[name].Scale(scale)
		if err != nil {
			logrus.Fatalf("Failed to set the scale %s=%d: %v", name, scale, err)
		}
	}
}

func wait() {
	<-make(chan interface{})
}
