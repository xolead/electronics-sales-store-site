package dbwork_test

import (
	"authoriz-service/pkg/dbwork"
	"context"
	"crypto/rand"
	"encoding/base64"
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
	"golang.org/x/crypto/bcrypt"
)

func TestMain(t *testing.T) {
	db, clean, err := setupTestDB()
	assert.NoError(t, err)
	defer clean()

	os.Setenv("expires_refresh", "10")
	EnableCheckSession_Success(t, db)
	CheckSession_NotFound(t, db)
	StopSession_Success(t, db)
	CreateRefresh_Success(t, db)
	CheckCollisionRefresh(t, db)
}

func EnableCheckSession_Success(t *testing.T, db *dbwork.DataBase) {
	GUID := uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.EnableSession(ctx, GUID)
	assert.NoError(t, err)

	err = db.CheckActiveSession(ctx, GUID)
	assert.NoError(t, err)
}

func CheckSession_NotFound(t *testing.T, db *dbwork.DataBase) {
	GUID := uuid.New().String()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.CheckActiveSession(ctx, GUID)
	assert.Error(t, err)
	assert.Equal(t, dbwork.SessionIsNotActive, err)
}

func StopSession_Success(t *testing.T, db *dbwork.DataBase) {
	GUID := uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.EnableSession(ctx, GUID)
	assert.NoError(t, err)

	err = db.StopSession(ctx, GUID)
	assert.NoError(t, err)

	err = db.StopSession(ctx, uuid.NewString())
	assert.NoError(t, err)

	err = db.EnableSession(ctx, GUID)
	assert.NoError(t, err)
}
func CreateRefresh_Success(t *testing.T, db *dbwork.DataBase) {
	GUID := uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var refresh [32]byte

	_, err := rand.Read(refresh[:])
	assert.NoError(t, err)

	strRefresh := base64.StdEncoding.EncodeToString(refresh[:])
	hashRefresh, err := bcrypt.GenerateFromPassword([]byte(strRefresh), bcrypt.DefaultCost)
	assert.NoError(t, err)

	err = db.CreateRefresh(ctx, string(hashRefresh), GUID)
	assert.NoError(t, err)
}

func CheckCollisionRefresh(t *testing.T, db *dbwork.DataBase) {
	GUID := uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var refresh [32]byte

	_, err := rand.Read(refresh[:])
	assert.NoError(t, err)

	strRefresh := base64.StdEncoding.EncodeToString(refresh[:])
	hashRefresh, err := bcrypt.GenerateFromPassword([]byte(strRefresh), bcrypt.DefaultCost)
	assert.NoError(t, err)

	err = db.CheckCollisionRefresh(ctx, string(hashRefresh))
	assert.NoError(t, err)

	err = db.CreateRefresh(ctx, string(hashRefresh), GUID)
	assert.NoError(t, err)

	err = db.CheckCollisionRefresh(ctx, string(hashRefresh))
	assert.Error(t, err)
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
