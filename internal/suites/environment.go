package suites

import (
	"fmt"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/authelia/authelia/v4/internal/utils"
)

//nolint:unparam
func waitUntilServiceLogDetected(
	interval time.Duration,
	timeout time.Duration,
	dockerEnvironment *DockerEnvironment,
	service string,
	logPatterns []string) error {
	log.Debug("Waiting for service " + service + " to be ready...")

	err := utils.CheckUntil(5*time.Second, 1*time.Minute, func() (bool, error) {
		logs, err := dockerEnvironment.Logs(service, []string{"--tail", "20"})
		fmt.Printf(".")

		if err != nil {
			return false, err
		}
		for _, pattern := range logPatterns {
			if strings.Contains(logs, pattern) {
				return true, nil
			}
		}
		return false, nil
	})

	fmt.Print("\n")

	return err
}

func waitUntilAutheliaBackendIsReady(dockerEnvironment *DockerEnvironment) error {
	return waitUntilServiceLogDetected(
		5*time.Second,
		90*time.Second,
		dockerEnvironment,
		"authelia-backend",
		[]string{"Initializing server for"})
}

func waitUntilAutheliaFrontendIsReady(dockerEnvironment *DockerEnvironment) error {
	return waitUntilServiceLogDetected(
		5*time.Second,
		90*time.Second,
		dockerEnvironment,
		"authelia-frontend",
		[]string{"dev server running at", "ready in"})
}

func waitUntilK3DIsReady(dockerEnvironment *DockerEnvironment) error {
	return waitUntilServiceLogDetected(
		5*time.Second,
		90*time.Second,
		dockerEnvironment,
		"k3d",
		[]string{"API listen on [::]:2376"})
}

func waitUntilSambaIsReady(dockerEnvironment *DockerEnvironment) error {
	return waitUntilServiceLogDetected(
		5*time.Second,
		90*time.Second,
		dockerEnvironment,
		"sambaldap",
		[]string{"samba entered RUNNING state"})
}

func waitUntilAutheliaIsReady(dockerEnvironment *DockerEnvironment, suite string) error {
	log.Info("Waiting for Authelia to be ready...")

	if err := waitUntilAutheliaBackendIsReady(dockerEnvironment); err != nil {
		return err
	}

	if os.Getenv("CI") != t && suite != "CLI" {
		if err := waitUntilAutheliaFrontendIsReady(dockerEnvironment); err != nil {
			return err
		}
	}

	if suite == "ActiveDirectory" {
		if err := waitUntilSambaIsReady(dockerEnvironment); err != nil {
			return err
		}
	}

	log.Info("Authelia is now ready!")

	return nil
}
