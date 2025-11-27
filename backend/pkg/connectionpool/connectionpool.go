package connectionpool

import (
	cloudstorage "electronic/pkg/cloud_storage"
	"electronic/pkg/dbwork"
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
