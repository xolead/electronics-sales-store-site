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
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"electronic/pkg/models"
)

type DataBase interface {
	CreateUser(login, password string) error
	VerifyPassword(login, password string) (bool, error)
	CreateProduct(pr models.Product) (int, error)
	UpdateProduct()
	DeleteProduct(id int) error
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

	postgres := &postgreSQL{db}

	return postgres, nil
}

func (postgres *postgreSQL) RunMigrations(path string) error {
	driver, err := postgre.WithInstance(postgres.DB, &postgre.Config{})
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

	_, err = postgres.Exec(createUserQuery, login, hashPassword)
	if err != nil {
		return err
	}

	return nil
}

func (postgres *postgreSQL) getUserID(login string) (int, error) {
	id := -1
	getUserQuery := `SELECT id FROM users WHERE login=$1`
	if err := postgres.QueryRow(getUserQuery, login).Scan(&id); err != nil {
		return -1, err
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
	                (name, description, parameters, count)
	                VALUES($1, $2, $3, $4) RETURNING id;`
	id := -1
	if err := postgres.QueryRow(createQueryProduct, pr.Name, pr.Description, pr.Parameters, pr.Count).Scan(&id); err != nil {
		return -1, err
	}
	createQueryImage := `INSERT INTO product_image
	                     (product_id, name, key)
	                     VALUES($1, $2, $3);`

	for _, image := range pr.Images {
		_, err := postgres.Exec(createQueryImage, id, image.Name, image.Key)
		if err != nil {
			_, err = postgres.Exec(createQueryImage, id, image.Name, image.Key)
			if err != nil {
				return -1, err
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
		return err
	}

	return nil

}

func (postgres *postgreSQL) DeleteProduct(id int) error {
	deleteQueryImage := `DELETE FROM product_image WHERE product_id=$1`
	if _, err := postgres.Exec(deleteQueryImage, id); err != nil {
		return err
	}

	deleteQueryProduct := `DELETE FROM product WHERE id=$1`
	if _, err := postgres.Exec(deleteQueryProduct, id); err != nil {
		return err
	}

	return nil
}

func (postgres *postgreSQL) ReadProduct(id int) (models.Product, error) {
	selectQueryProduct := `SELECT id, name, description, parameters, count FROM product WHERE id=$1`
	pr := models.Product{}

	if err := postgres.QueryRow(selectQueryProduct, id).Scan(&pr.ID, &pr.Name, &pr.Description, &pr.Parameters, &pr.Count); err != nil {
		return pr, err
	}

	selectQueryImages := `SELECT id, product_id, name, key FROM product_image WHERE product_id=$1`

	rows, err := postgres.Query(selectQueryImages, id)
	if err != nil {
		return pr, err
	}
	defer rows.Close()

	for rows.Next() {
		image := models.ProductImage{}
		if err = rows.Scan(&image.ID, &image.ProductID, &image.Name, &image.Key); err != nil {
			return pr, err
		}

		pr.Images = append(pr.Images, image)

	}
	return pr, nil
}

func (postgres *postgreSQL) ReadListProduct() ([]models.Product, error) {
	selectQueryProduct := `SELECT id, name, description, parameters, count FROM product`
	pr := make([]models.Product, 0)

	rows, err := postgres.Query(selectQueryProduct)
	if err != nil {
		return pr, err
	}
	defer rows.Close()
	for rows.Next() {
		product := models.Product{}
		if err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Parameters, &product.Count); err != nil {
			return pr, err
		}
		selectQueryImages := `SELECT id, product_id, name, key FROM product_image WHERE product_id=$1`

		rowsImages, err := postgres.Query(selectQueryImages, product.ID)
		if err != nil {
			return pr, err
		}
		defer rowsImages.Close()
		for rowsImages.Next() {
			image := models.ProductImage{}
			if err = rowsImages.Scan(&image.ID, &image.ProductID, &image.Name, &image.Key); err != nil {
				return pr, err
			}
			product.Images = append(product.Images, image)
		}
		pr = append(pr, product)
	}

	return pr, nil
}
