package suites

import (
	"fmt"
	"time"
)

var bypassAllSuiteName = "BypassAll"

func init() {
	dockerEnvironment := NewDockerEnvironment([]string{
		"docker-compose.yml",
		"internal/suites/BypassAll/docker-compose.yml",
		"example/compose/authelia/docker-compose.backend.{}.yml",
		"example/compose/authelia/docker-compose.frontend.{}.yml",
		"example/compose/nginx/backend/docker-compose.yml",
		"example/compose/nginx/portal/docker-compose.yml",
		"example/compose/httpbin/docker-compose.yml",
		"example/compose/smtp/docker-compose.yml",
		"example/compose/duo-api/docker-compose.yml",
	})

	setup := func(suitePath string) error {
		if err := dockerEnvironment.Up(); err != nil {
			return err
		}

		return waitUntilAutheliaBackendIsReady(dockerEnvironment)
	}

	onSetupTimeout := func() error {
		backendLogs, err := dockerEnvironment.Logs("authelia-backend", nil)
		if err != nil {
			return err
		}
		fmt.Println(backendLogs)

		frontendLogs, err := dockerEnvironment.Logs("authelia-frontend", nil)
		if err != nil {
			return err
		}
		fmt.Println(frontendLogs)
		return nil
	}

	teardown := func(suitePath string) error {
		return dockerEnvironment.Down()
	}

	GlobalRegistry.Register(bypassAllSuiteName, Suite{
		SetUp:           setup,
		SetUpTimeout:    5 * time.Minute,
		OnSetupTimeout:  onSetupTimeout,
		TestTimeout:     1 * time.Minute,
		TearDown:        teardown,
		TearDownTimeout: 2 * time.Minute,
		Description:     "This suite has been created to test Authelia with a bypass policy on all resources",
	})
}
