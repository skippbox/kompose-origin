package app

import (
    "bytes"
    "io/ioutil"
    "os"
    "strings"
    "text/template"
    "regexp"

    "github.com/Sirupsen/logrus"
)

/* Ancilliary helper functions to interface with the commands interface */

/**
 * Retrieve the kubernetes server based on .kuberconfig or use a default of
 * 127.0.0.1:8080.
 */
func getK8sServer(file string) string {
    if len(file) == 0 {
        file = ".kuberconfig"
    }

    roughData,err := ioutil.ReadFile(file)

    if err != nil {
        logrus.Debugf("Cannot read server data, defaulting to: 127.0.0.1:8080")
        return "127.0.0.1:8080"
    }

    server := strings.TrimSpace(string(roughData))
    foundPort, err := regexp.MatchString(".+:[\\d]+", server)
    if !foundPort || err != nil {
        server += ":8080"
    }

    return server
}

/**
 * Generate Helm Chart configuration
 */
func generateHelm(filename string, svcname string) error {
    type ChartDetails struct {
        Name string
    }

    dirName := strings.Replace(filename, ".yml", "", 1)
    details := ChartDetails{dirName}
    manifestDir := dirName + string(os.PathSeparator) + "manifests"
    dir,err := os.Open(dirName)

    /* Setup the initial directories/files */
    if err == nil {
        _ = dir.Close()
    }

    if err != nil {
        err = os.Mkdir(dirName, 0755)
        if err != nil {
            return err
        }

        err = os.Mkdir(manifestDir, 0755)
        if err != nil {
            return err
        }

        /* Create the readme file */
        readme := "This chart was created by Kompose\n"
        err = ioutil.WriteFile(dirName + string(os.PathSeparator) + "README.md", []byte(readme), 0644)
        if err != nil {
            return err
        }

        /* Create the Chart.yaml file */
        chart := `name: {{.Name}}
description: A generated Helm Chart from Skippbox Kompose
version: 0.0.1
source:
home:
`

        t, err := template.New("ChartTmpl").Parse(chart)
        if err != nil {
            logrus.Fatalf("Failed to generate Chart.yaml template: %s\n", err)
        }
        var chartData bytes.Buffer
        _ = t.Execute(&chartData, details)

        err = ioutil.WriteFile(dirName + string(os.PathSeparator) + "Chart.yaml", chartData.Bytes(), 0644)
        if err != nil {
            return err
        }
    }

    /* Copy all yaml files into the newly created manifests directory */
    infile, err := ioutil.ReadFile(svcname + "-rc.yaml")
    if err != nil {
        logrus.Infof("Error reading %s: %s\n", svcname + "-rc.yaml", err)
        return err
    }

    err = ioutil.WriteFile(manifestDir + string(os.PathSeparator) + svcname + "-rc.yaml", infile, 0644)
    if err != nil {
        return err
    }

    /* The svc file is optional */
    infile, err = ioutil.ReadFile(svcname + "-svc.yaml")
    if err == nil {
        err = ioutil.WriteFile(manifestDir + string(os.PathSeparator) + svcname + "-svc.yaml", infile, 0644)
        if err != nil {
           return err
        }
    }

    return nil
}