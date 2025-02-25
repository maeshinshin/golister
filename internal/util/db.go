package util

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DBInfo struct {
	DB_DATABASE string
	DB_USERNAME string
	DB_PASSWORD string
	Db_PORT     string
	Db_HOST     string
}

func MustStartMySQLContainer(DbInfo *DBInfo) (func(context.Context, ...testcontainers.TerminateOption) error, error) {

	dbContainer, err := mysql.Run(context.Background(),
		"mysql:8.0.36",
		mysql.WithDatabase(DbInfo.DB_DATABASE),
		mysql.WithUsername(DbInfo.DB_USERNAME),
		mysql.WithPassword(DbInfo.DB_PASSWORD),
		testcontainers.WithWaitStrategy(wait.ForLog("port: 3306  MySQL Community Server - GPL").WithStartupTimeout(30*time.Second)),
	)

	if err != nil {
		return nil, err
	}

	DbInfo.Db_HOST, err = dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "3306/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}
	DbInfo.Db_PORT = dbPort.Port()

	return dbContainer.Terminate, err
}
