package testutil

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ContainerData struct {
	host string
	port string
}

func (c *ContainerData) Host() string {
	return c.host
}

func (c *ContainerData) Port() string {
	return c.port
}

func MustStartMySQLContainer(dbName, dbPwd, dbUser string) (func(context.Context, ...testcontainers.TerminateOption) error, *ContainerData, error) {

	dbContainer, err := mysql.Run(context.Background(),
		"mysql:8.0.36",
		mysql.WithDatabase(dbName),
		mysql.WithUsername(dbUser),
		mysql.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(wait.ForLog("port: 3306  MySQL Community Server - GPL").WithStartupTimeout(30*time.Second)),
	)

	if err != nil {
		return nil, nil, err
	}

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, nil, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "3306/tcp")
	if err != nil {
		return dbContainer.Terminate, nil, err
	}

	return dbContainer.Terminate, &ContainerData{dbHost, dbPort.Port()}, err
}
