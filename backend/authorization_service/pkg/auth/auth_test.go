package auth_test

import (
	"authoriz-service/pkg/auth"
	"authoriz-service/pkg/dbwork"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestMain(t *testing.T) {
	db, clean, err := setupTestDB()
	assert.NoError(t, err)
	defer clean()

	os.Setenv("expires_jwt", "2")
	os.Setenv("jwt_secret", "MiNeLoshadi")
	os.Setenv("expires_refresh", "10")

	CreateAccessToken_Success(t, db)
	CreateAccessToken_Invalid(t, db)
	CreateRefreshToken_Success(t, db)
}

func CreateAccessToken_Success(t *testing.T, db *dbwork.DataBase) {
	GUID := uuid.NewString()
	access, err := auth.CreateAccessToken(db, GUID, false)
	assert.NoError(t, err)

	NewGUID, admin, err := auth.CheckAccessToken(access)
	assert.NoError(t, err)
	assert.Equal(t, GUID, NewGUID)
	assert.Equal(t, false, admin)

}

func CreateAccessToken_Invalid(t *testing.T, db *dbwork.DataBase) {

	GUID, admin, err := auth.CheckAccessToken(uuid.NewString())
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil.String(), GUID)
	assert.Equal(t, false, admin)
}

func CreateRefreshToken_Success(t *testing.T, db *dbwork.DataBase) {
	GUID := uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	refresh, err := auth.CreateRefreshToken(ctx, db, GUID)
	assert.NoError(t, err)

	err = db.CheckRefreshToken(ctx, GUID, refresh)
	assert.NoError(t, err)
}

func setupTestDB() (*dbwork.DataBase, func(), error) {
	dbName := "testdb"
	dbUser := "test"
	dbPassword := "pass"
	ctx := context.Background()
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:15-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, nil, err
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, nil, err
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, err
	}

	db, err := dbwork.NewPostgreSQL(ctx, dbwork.PostgreSQLConfig{
		User:     dbUser,
		Password: dbPassword,
		Host:     host,
		Port:     port.Int(),
		DBName:   dbName,
		SSLMode:  "disable",
	})

	if err != nil {
		return nil, nil, err
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, nil, fmt.Errorf("cannot get current file path")
	}

	dir := filepath.Dir(filename)
	migrationsPath := filepath.Join(dir, "migrations")

	err = db.RunMigrations(ctx, migrationsPath)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		db.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}

	return db, cleanup, nil

}
