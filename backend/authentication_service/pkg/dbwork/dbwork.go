package dbwork

import (
	"auth-service/pkg/models"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var (
	LoginBusy            = errors.New("Данный логин уже занят")
	LoginNotFound        = errors.New("Неправильный логин или пароль")
	PasswordIsNotCorrect = errors.New("Неправильный логин или пароль")
)

func init() {
	log.Logger = log.With().Str("package", "dbwork").Logger()
}

type PostgreSQLConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
	SSLMode  string
}

type DataBase struct {
	pool *pgxpool.Pool
}

func LoadPSQLConfig() PostgreSQLConfig {
	port, err := strconv.Atoi(os.Getenv("port"))
	if err != nil {
		log.Fatal().Msgf("LoadPSQLConfig getPort: %v", err)
	}
	return PostgreSQLConfig{
		User:     os.Getenv("user"),
		Password: os.Getenv("password"),
		Host:     os.Getenv("host"),
		Port:     port,
		DBName:   os.Getenv("dbname"),
		SSLMode:  os.Getenv("sslmode"),
	}
}

func NewPostgreSQL(ctx context.Context, config PostgreSQLConfig) (*DataBase, error) {
	log.Info().Msgf("Post cfg: %v", config)
	connStr := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("dbwork/NewPostgreSQL poolConfig: %v", err)
	}

	poolConfig.MaxConnIdleTime = time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("dbwork/NewPostgreSQL newPool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("dbwork/NewPostgreSQL ping: %v", err)
	}
	log.Info().Msg("Подключение к базе данных успешно!")

	return &DataBase{pool: pool}, nil
}

func (db *DataBase) Close() {
	if db.pool != nil {
		db.pool.Close()
		log.Info().Msg("Соеднинение с базой данных закрыто")
	}
}

func (db *DataBase) RunMigrations(ctx context.Context, path string) error {
	sqlDB := stdlib.OpenDBFromPool(db.pool)
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("Ошибка установки диалекта: %w", err)
	}

	if path == "" {
		path = "pkg/dbwork/migrations"
	}

	if err := goose.Up(sqlDB, path); err != nil {
		return fmt.Errorf("Ошибка применения миграций: %w", err)
	}

	log.Info().Msg("Миграции успешно применены")
	return nil
}

func (db *DataBase) CreateUser(ctx context.Context, login, password string) (uuid.UUID, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, fmt.Errorf("dbwork/CreateUser generateHashPassword: %v", err)
	}

	createQuery := `INSERT INTO users
	                       (login, password, registration_date)
							VALUES ($1, $2, $3)
							RETURNING id`
	date := time.Now()
	id := uuid.Nil

	err = db.pool.QueryRow(ctx, createQuery, login, string(hashPassword), date).Scan(&id)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return uuid.Nil, LoginBusy
		}
		return uuid.Nil, fmt.Errorf("dbwork/CreateUser QueryRow: %v", err)
	}

	return id, nil

}

func (db *DataBase) VerifyPassword(ctx context.Context, login, password string) (uuid.UUID, bool, error) {

	selectQuery := `SELECT id, login, password, admin
	                FROM users
					WHERE login = $1`

	user := models.User{}

	err := db.pool.QueryRow(ctx, selectQuery, login).Scan(&user.ID, &user.Login, &user.Password, &user.Admin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, false, LoginNotFound
		}
		return uuid.Nil, false, fmt.Errorf("dbwork/VerifyPassword selectQuery: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err == nil {
		return user.ID, user.Admin, nil
	}

	return uuid.Nil, false, PasswordIsNotCorrect
}

func (db *DataBase) MakeAdmin(ctx context.Context, login string) (uuid.UUID, error) {
	updateQuery := `UPDATE users
	                       SET admin = true
						   WHERE login=$1
						   RETURNING id`
	id := uuid.Nil

	err := db.pool.QueryRow(ctx, updateQuery, login).Scan(&id)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
