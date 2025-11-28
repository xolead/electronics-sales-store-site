package connectionpool

import (
	cloudstorage "core-service/pkg/cloud_storage"
	"core-service/pkg/dbwork"
)

var (
	dataBase dbwork.DataBase
	s3s      cloudstorage.CloudStorage
)

func init() {
	var err error
	dataBase, err = dbwork.NewPostgreSQL(dbwork.LoadPSQLConfig())
	if err != nil {
		panic(err)
	}

	dataBase.RunMigrations("")

	s3s, err = cloudstorage.NewS3S(cloudstorage.LoadConfig())
	if err != nil {
		panic(err)
	}
}

type ConnectionPool struct{}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{}
}

func (cp *ConnectionPool) GetDataBase() dbwork.DataBase {
	return dataBase
}

func (cp *ConnectionPool) GetS3Storage() cloudstorage.CloudStorage {
	return s3s
}

func (cp *ConnectionPool) Close() {
	dataBase.Close()
}
