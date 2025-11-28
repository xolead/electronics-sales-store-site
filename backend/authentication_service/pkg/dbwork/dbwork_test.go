package dbwork_test

import (
	"auth-service/pkg/dbwork"
	"context"
	"fmt"
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

	CreateUser_Success(t, db)
	CreateUser_DuplicateLogin(t, db)
	VerifyPassword_Success(t, db)
	VerifyPassword_LoginNotFound(t, db)
	VerifyPassword_PasswordIsNotCorrect(t, db)
	MakeAdmin_Success(t, db)

}

func CreateUser_Success(t *testing.T, db *dbwork.DataBase) {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	id, err := db.CreateUser(ctx, "login123", "pass88888")
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func CreateUser_DuplicateLogin(t *testing.T, db *dbwork.DataBase) {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	id, err := db.CreateUser(ctx, "loginlogin", "passs123123")
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	id, err = db.CreateUser(ctx, "loginlogin", "pas7777")
	assert.Error(t, err)
	assert.Equal(t, dbwork.LoginBusy, err)
	assert.Equal(t, uuid.Nil, id)
}

func VerifyPassword_Success(t *testing.T, db *dbwork.DataBase) {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	id, err := db.CreateUser(ctx, "-1wx???--LP", "-_-_-_-_-")
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	id, admin, err := db.VerifyPassword(ctx, "-1wx???--LP", "-_-_-_-_-")
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
	assert.Equal(t, false, admin)
}

func VerifyPassword_LoginNotFound(t *testing.T, db *dbwork.DataBase) {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	id, admin, err := db.VerifyPassword(ctx, "abobaaboba", "-------------")
	assert.Error(t, err)
	assert.Equal(t, dbwork.LoginNotFound, err)
	assert.Equal(t, uuid.Nil, id)
	assert.Equal(t, false, admin)
}

func VerifyPassword_PasswordIsNotCorrect(t *testing.T, db *dbwork.DataBase) {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	id, err := db.CreateUser(ctx, "Za_Alians", "Za_Ordy")
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	id, admin, err := db.VerifyPassword(ctx, "Za_Alians", "Orda_lox")
	assert.Error(t, err)
	assert.Equal(t, dbwork.PasswordIsNotCorrect, err)
	assert.Equal(t, uuid.Nil, id)
	assert.Equal(t, false, admin)
}

func MakeAdmin_Success(t *testing.T, db *dbwork.DataBase) {
	ctx, cancel := context.WithTimeout(context.TODO(), 4*time.Second)
	defer cancel()

	id, err := db.CreateUser(ctx, "LOLKEK", "Chebyrek")
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)

	newID, err := db.MakeAdmin(ctx, "LOLKEK")
	assert.NoError(t, err)
	assert.Equal(t, id, newID)

	id, admin, err := db.VerifyPassword(ctx, "LOLKEK", "Chebyrek")
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
	assert.Equal(t, true, admin)
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
