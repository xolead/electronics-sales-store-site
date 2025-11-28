package dbwork

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

var (
	DuplicateRefresh    = errors.New("Данный refrsh токен занят")
	SessionIsNotActive  = errors.New("Сессия не активна")
	RefreshIsNotActive  = errors.New("Refresh токен неактивен")
	InvalidRefreshToken = errors.New("Неправильный refredh токе")
)

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

func (db *DataBase) CheckCollisionRefresh(ctx context.Context, hashRefresh string) error {
	selectQuery := `SELECT id
	                       FROM refresh
						   WHERE refresh_token=$1`
	id := -1
	if err := db.pool.QueryRow(ctx, selectQuery, hashRefresh).Scan(&id); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return DuplicateRefresh
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("dbwork/CheckCollisionRefresh QueryRow: %v", err)
	}

	if id == -1 {
		return nil
	}
	return DuplicateRefresh
}

func (db *DataBase) CreateRefresh(ctx context.Context, hashRefresh, GUID string) error {

	if err := db.StopWorkerRefresh(ctx, GUID); err != nil {
		return fmt.Errorf("CreateRefresh: %v", err)
	}

	createQuery := `INSERT INTO refresh
	                       (user_id, refresh_token, worker, expires_at)
	                       VALUES($1, $2, TRUE, $3);`

	expires, err := strconv.Atoi(os.Getenv("expires_refresh"))
	if err != nil {
		return fmt.Errorf("dbwork/CreateRefresh os.Getenv: %v", err)
	}
	_, err = db.pool.Exec(ctx, createQuery, GUID,
		hashRefresh,
		time.Now().Add(time.Duration(expires)*time.Hour))

	if err != nil {
		return fmt.Errorf("dbwork/CreateRefresh Exec: %v", err)
	}

	return nil
}

func (db *DataBase) StopWorkerRefresh(ctx context.Context, GUID string) error {
	updateQuery := `UPDATE refresh
	                       SET worker=false
					       WHERE user_id = $1 AND worker=true`

	if _, err := db.pool.Exec(ctx, updateQuery, GUID); err != nil {
		return fmt.Errorf("dbwork/StopWorkerRefresh exec: %v", err)
	}

	return nil
}

func (db *DataBase) CheckRefreshToken(ctx context.Context, GUID, refresh string) error {
	selectQuery := `SELECT refresh_token
	                FROM refresh
	                WHERE user_id=$1 AND worker=TRUE`
	hashRefresh := ""

	if err := db.pool.QueryRow(ctx, selectQuery, GUID).Scan(&hashRefresh); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return RefreshIsNotActive
		}
		return fmt.Errorf("dbwork/CheckRefreshToken QueryRow: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashRefresh), []byte(refresh)); err != nil {
		return InvalidRefreshToken
	}

	return nil

}

func (db *DataBase) EnableSession(ctx context.Context, GUID string) error {
	updateQuery := `UPDATE session
	                SET active=TRUE
	                WHERE user_id=$1`

	result, err := db.pool.Exec(ctx, updateQuery, GUID)
	if err != nil {
		return fmt.Errorf("dbwork/EnableSession Exec: %v", err)
	}

	if result.RowsAffected() > 0 {
		return nil
	}

	if err = db.CreateSession(ctx, GUID); err != nil {
		return fmt.Errorf("EnableSession: %v", err)
	}

	return nil
}

func (db *DataBase) CreateSession(ctx context.Context, GUID string) error {
	createQuery := `INSERT INTO session
	                       (user_id, active)
						   VALUES($1, TRUE)`

	if _, err := db.pool.Exec(ctx, createQuery, GUID); err != nil {
		return fmt.Errorf("dbwork/CreateSession Exec: %v", err)
	}

	return nil
}

func (db *DataBase) CheckActiveSession(ctx context.Context, GUID string) error {
	selectQuery := `SELECT id
	                       FROM session
	                       WHERE user_id=$1 AND active=TRUE`
	id := -1

	if err := db.pool.QueryRow(ctx, selectQuery, GUID).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return SessionIsNotActive
		}
		return fmt.Errorf("dbwork/CheckActiveSession QueryRow: %v", err)
	}

	return nil
}

func (db *DataBase) StopSession(ctx context.Context, GUID string) error {
	updateQuery := `UPDATE session
	                SET active = FALSE
	                WHERE user_id = $1`
	if _, err := db.pool.Exec(ctx, updateQuery, GUID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("dbwork/StopSession Exec: %v", err)
	}

	return nil

}
