package dbwork

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	postgre "github.com/golang-migrate/migrate/v4/database/postgres"
	"golang.org/x/crypto/bcrypt"
)

type DataBase interface {
	CreateUser(login, password string) error
	CreateProduct()
	UpdateProduct()
	DeleteProduct()
	ReadProduct()
	ReadListProduct()
	RunMigrations()
	Close()
}

type postgreSQL struct {
	db *sql.DB
}

type PostgreSQLConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
	SSLMode  string
}

func NewPostgreSQL(config PostgreSQLConfig) (*DataBase, error) {
	connStr := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	postgres := &postgreSQL{db: db}

	return postgres, nil
}

func (postgres *postgreSQL) RunMigrations(path string) error {
	driver, err := postgre.WithInstance(postgres.db, &postgre.Config{})
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if path == "" {
		path = "file://" + filepath.Join(wd, "pkg/dbwork/migrations")
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return err
	}

	if dirty {
		if err := m.Force(int(version)); err != nil {
			return err
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		if dirtyErr, ok := err.(migrate.ErrDirty); ok {
			if err := m.Force(int(dirtyErr.Version)); err != nil {
				return err
			}
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				return err
			}
		}
		return err
	}
	return nil
}

func (postgres *postgreSQL) Close() {
	postgres.db.Close()
}

func (postgres *postgreSQL) CreateUser(login, password string) error {
	createUserQuery := `INSERT INTO users
											(login, password)
											VALUES($1, $2);`
	id, err := postgres.getUserID(login)
	if err != nil {
		return err
	}

	if id != -1 {
		return fmt.Errorf("Пользователь с таким логином уже существует")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = postgres.db.Exec(createUserQuery, login, hashPassword)
	if err != nil {
		return err
	}

	return nil
}

func (postgres *postgreSQL) getUserID(login string) (int, error) {
	id := -1
	getUserQuery := `SELECT id FROM users WHERE login=$1`
	if err := postgres.db.QueryRow(getUserQuery, login).Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (postgres *postgreSQL) VerifyPassword(login, password string) (bool, error) {
	getUserQuery := `SELECT password FROM users WHERE login=$1`

	var realPassword []byte

	if err := postgres.db.QueryRow(getUserQuery, login).Scan(&realPassword); err != nil {
		return false, err
	}

	err := bcrypt.CompareHashAndPassword(realPassword, []byte(password))

	return err == nil, nil
}

func (postgres *postgreSQL) CreateProduct() {}

func (postgres *postgreSQL) UpdateProduct() {}

func (postgres *postgreSQL) DeleteProduct() {}

func (postgres *postgreSQL) ReadProduct() {}

func (postgres *postgreSQL) ReadListProduct() {}
