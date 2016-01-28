package app

import (
    "io/ioutil"
    "strings"
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