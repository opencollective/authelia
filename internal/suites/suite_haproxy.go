package suites

import (
	"fmt"
	"time"
)

var haproxySuiteName = "HAProxy"

func init() {
	dockerEnvironment := NewDockerEnvironment([]string{
		"docker-compose.yml",
		"internal/suites/HAProxy/docker-compose.yml",
		"example/compose/authelia/docker-compose.backend.{}.yml",
		"example/compose/authelia/docker-compose.frontend.{}.yml",
		"example/compose/nginx/backend/docker-compose.yml",
		"example/compose/haproxy/docker-compose.yml",
		"example/compose/smtp/docker-compose.yml",
	})

	setup := func(suitePath string) error {
		err := dockerEnvironment.Up()

		if err != nil {
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
		err := dockerEnvironment.Down()
		return err
	}

	GlobalRegistry.Register(haproxySuiteName, Suite{
		SetUp:           setup,
		SetUpTimeout:    5 * time.Minute,
		OnSetupTimeout:  onSetupTimeout,
		TestTimeout:     2 * time.Minute,
		TearDown:        teardown,
		TearDownTimeout: 2 * time.Minute,
	})
}
