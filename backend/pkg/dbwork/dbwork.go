package dbwork

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	postgre "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"electronic/pkg/models"
)

type DataBase interface {
	CreateUser(login, password string) error
	VerifyPassword(login, password string) (bool, error)
	CreateProduct(pr models.Product) (int, error)
	UpdateProduct()
	DeleteProduct(id int) ([]string, error)
	ReadProduct(id int) (models.Product, error)
	ReadListProduct() ([]models.Product, error)
	RunMigrations(path string) error
	ChangeCountProduct(id, changeCount int) error
	Close()
}

type postgreSQL struct {
	*sql.DB
}

type PostgreSQLConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
	SSLMode  string
}

func LoadPSQLConfig() PostgreSQLConfig {
	port, err := strconv.Atoi(os.Getenv("port"))
	if err != nil {
		fmt.Println(port)
		log.Fatal(err)
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

func (postgres *postgreSQL) Close() {
	postgres.DB.Close()
}

func NewPostgreSQL(config PostgreSQLConfig) (DataBase, error) {
	log.Println("Post cfg: ", config)
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
		return nil, fmt.Errorf("Ошибка соединения с базой данных: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Ошибка ping базы данных: %w", err)
	}

	postgres := &postgreSQL{db}

	err = postgres.RunMigrations("")
	if err != nil {
		return nil, fmt.Errorf("NewPostgreSQL ошибка миграции: %w", err)
	}

	return postgres, nil
}

func (postgres *postgreSQL) RunMigrations(path string) error {
	driver, err := postgre.WithInstance(postgres.DB, &postgre.Config{})
	if err != nil {
		return fmt.Errorf("Оишбка в создании драйвера для миграций: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Ошибка в Getwd, RunMigrations: %w", err)
	}

	if path == "" {
		path = "file://" + filepath.Join(wd, "pkg/dbwork/migrations")
	}

	m, err := migrate.NewWithDatabaseInstance(
		path,
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("Ошибка миграции NewWithDatabaseInstance: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("Ошибка в версии миграции, dirty: %w", err)
	}

	if dirty {
		if err := m.Force(int(version)); err != nil {
			return fmt.Errorf("Невозможность исправить dirty миграции: %w", err)
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		if dirtyErr, ok := err.(migrate.ErrDirty); ok {
			if err := m.Force(int(dirtyErr.Version)); err != nil {
				return fmt.Errorf("Dirty миграций: %w", err)
			}
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				return fmt.Errorf("Невозможность поднять миграцию: %w", err)
			}
		}
		return err
	}
	log.Println("Миграция прошла успешно")
	return nil
}

func (postgres *postgreSQL) CreateUser(login, password string) error {
	createUserQuery := `INSERT INTO users
											(login, password)
											VALUES($1, $2);`
	id, err := postgres.getUserID(login)
	if err != nil {
		return fmt.Errorf("CreateUser ошибка в getUserID: %w", err)
	}

	if id != -1 {
		return fmt.Errorf("Пользователь с таким логином уже существует")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Ошибка генерации хэша пароля: %w", err)
	}

	_, err = postgres.Exec(createUserQuery, login, hashPassword)
	if err != nil {
		return fmt.Errorf("Ошибка создани пользователя Exec: %w", err)
	}

	return nil
}

func (postgres *postgreSQL) getUserID(login string) (int, error) {
	id := -1
	getUserQuery := `SELECT id FROM users WHERE login=$1`
	if err := postgres.QueryRow(getUserQuery, login).Scan(&id); err != nil {
		return -1, fmt.Errorf("getUserID ошибка queryRow: %w", err)
	}
	return id, nil
}

func (postgres *postgreSQL) VerifyPassword(login, password string) (bool, error) {
	getUserQuery := `SELECT password FROM users WHERE login=$1`

	var realPassword []byte

	if err := postgres.QueryRow(getUserQuery, login).Scan(&realPassword); err != nil {
		return false, err
	}

	err := bcrypt.CompareHashAndPassword(realPassword, []byte(password))

	return err == nil, nil
}

func (postgres *postgreSQL) CreateProduct(pr models.Product) (int, error) {
	createQueryProduct := `INSERT INTO product
	                (name, description, parameters, count, price)
	                VALUES($1, $2, $3, $4, $5) RETURNING id;`
	id := -1
	if err := postgres.QueryRow(createQueryProduct, pr.Name, pr.Description, pr.Parameters, pr.Count, pr.Price).Scan(&id); err != nil {
		return -1, fmt.Errorf("CreateProduct ошибка QueryRow: %w", err)
	}
	createQueryImage := `INSERT INTO product_image
	                     (product_id, name, key)
	                     VALUES($1, $2, $3);`

	for _, image := range pr.Images {
		_, err := postgres.Exec(createQueryImage, id, image.Name, image.Key)
		if err != nil {
			_, err = postgres.Exec(createQueryImage, id, image.Name, image.Key)
			if err != nil {
				return -1, fmt.Errorf("CreateProduct ошибка добавления изображений: %w", err)
			}
		}
	}
	return id, nil
}

func (postgres *postgreSQL) UpdateProduct() {}

func (postgres *postgreSQL) ChangeCountProduct(id, countChange int) error {
	updateQueryCount := `UPDATE product SET count = count + $1 WHERE id = $2`
	if countChange == 0 {
		return nil
	}
	if _, err := postgres.Exec(updateQueryCount, countChange, id); err != nil {
		return fmt.Errorf("UpdateProduct ошибка exec: %w", err)
	}

	return nil

}

func (postgres *postgreSQL) DeleteProduct(id int) ([]string, error) {
	keys := make([]string, 0)

	selectQuery := `SELECT key FROM product_image WHERE product_id=$1`

	rows, err := postgres.Query(selectQuery, id)
	if err != nil {
		return keys, fmt.Errorf("DeleteProduct ошибка query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		key := ""
		if err = rows.Scan(&key); err != nil {
			return keys, fmt.Errorf("DeleteProduct ошибка scan product: %w", err)
		}
		keys = append(keys, key)
	}

	deleteQueryImage := `DELETE FROM product_image WHERE product_id=$1`
	if _, err := postgres.Exec(deleteQueryImage, id); err != nil {
		return keys, fmt.Errorf("DeleteProduct ошибка exec удаления image: %w", err)
	}

	deleteQueryProduct := `DELETE FROM product WHERE id=$1`
	if _, err := postgres.Exec(deleteQueryProduct, id); err != nil {
		return keys, fmt.Errorf("DeleteProduct ошибка exec удаления product: %w", err)
	}

	return keys, nil
}

func (postgres *postgreSQL) ReadProduct(id int) (models.Product, error) {
	selectQueryProduct := `SELECT id, name, description, parameters, count, price FROM product WHERE id=$1`
	pr := models.Product{}

	if err := postgres.QueryRow(selectQueryProduct, id).Scan(&pr.ID, &pr.Name, &pr.Description, &pr.Parameters, &pr.Count, &pr.Price); err != nil {
		return pr, fmt.Errorf("ReadProduct ошибка queryrow product: %w", err)
	}

	selectQueryImages := `SELECT id, product_id, name, key FROM product_image WHERE product_id=$1`

	rows, err := postgres.Query(selectQueryImages, id)
	if err != nil {
		return pr, fmt.Errorf("ReadProduct ошибка query image: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		image := models.ProductImage{}
		if err = rows.Scan(&image.ID, &image.ProductID, &image.Name, &image.Key); err != nil {
			return pr, fmt.Errorf("ReadProduct ошибка scan image: %w", err)
		}

		pr.Images = append(pr.Images, image)

	}
	return pr, nil
}

func (postgres *postgreSQL) ReadListProduct() ([]models.Product, error) {
	selectQueryProduct := `SELECT id, name, description, parameters, count, price FROM product`
	pr := make([]models.Product, 0)

	rows, err := postgres.Query(selectQueryProduct)
	if err != nil {
		return pr, fmt.Errorf("ReadListProduct ошибка Query product: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		product := models.Product{}
		if err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Parameters, &product.Count, &product.Price); err != nil {
			return pr, fmt.Errorf("ReadListProduct ошибка scan product: %w", err)
		}
		selectQueryImages := `SELECT id, product_id, name, key FROM product_image WHERE product_id=$1`

		rowsImages, err := postgres.Query(selectQueryImages, product.ID)
		if err != nil {
			return pr, fmt.Errorf("ReadListProduct ошибка query image: %w", err)
		}
		defer rowsImages.Close()
		for rowsImages.Next() {
			image := models.ProductImage{}
			if err = rowsImages.Scan(&image.ID, &image.ProductID, &image.Name, &image.Key); err != nil {
				return pr, fmt.Errorf("ReadListProduct ошибка scan image: %w", err)
			}
			product.Images = append(product.Images, image)
		}
		pr = append(pr, product)
	}

	return pr, nil
}
