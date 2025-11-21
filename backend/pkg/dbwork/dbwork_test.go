package dbwork_test

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"electronic/pkg/dbwork"
	"electronic/pkg/models"
)

func setupTestDB(ctx context.Context) (dbwork.DataBase, func(), error) {
	dbName := "testdb"
	dbUser := "test"
	dbPassword := "pass"
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
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

	config := dbwork.PostgreSQLConfig{
		User:     dbUser,
		Password: dbPassword,
		Host:     host,
		Port:     port.Int(),
		DBName:   dbName,
		SSLMode:  "disable",
	}

	db, err := dbwork.NewPostgreSQL(config)
	if err != nil {
		return nil, nil, err
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, nil, fmt.Errorf("cannot get current file path")
	}

	dir := filepath.Dir(filename)
	migrationsPath := "file://" + filepath.Join(dir, "migrations")
	err = db.RunMigrations(migrationsPath)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		db.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatal(err)
		}
	}

	return db, cleanup, nil
}

func TestPostgreSQL(t *testing.T) {
	ctx := context.Background()
	db, cleanup, err := setupTestDB(ctx)
	assert.NoError(t, err)
	defer cleanup()

	db.RunMigrations("")

	testImage := models.ProductImage{
		Name: "lol.txt",
		Key:  uuid.New().String(),
	}

	testProduct := models.Product{
		Name:        "Айфон",
		Description: "Чёткий",
		Parameters:  "12 дюймов, черный",
	}
	testProduct.Images = append(testProduct.Images, testImage)

	CreateAndReadProduct(t, testProduct, db)

	testProduct2 := models.Product{
		Name:        "Рено логан",
		Description: "Не бит не крашен",
		Parameters:  "Чёрный",
	}

	testImage2 := models.ProductImage{
		Name: "xaxaxaxaxa9817hd.jpg",
		Key:  uuid.New().String(),
	}

	testImage3 := models.ProductImage{
		Name: "8137erfchwiudc.png",
		Key:  uuid.New().String(),
	}

	testProduct2.Images = append(testProduct2.Images, testImage2, testImage3)

	CreateAndReadProduct(t, testProduct2, db)

	testProducts := make([]models.Product, 0)
	testProducts = append(testProducts, testProduct, testProduct2)
	ReadAllProduct(t, testProducts, db)
	DeleteAndRead(t, db)
	ChangeCountProduct(t, db)

}

func CreateAndReadProduct(t *testing.T, TestProduct models.Product, db dbwork.DataBase) {
	id, err := db.CreateProduct(TestProduct)
	assert.NoError(t, err)
	DBProduct, err := db.ReadProduct(id)
	assert.NoError(t, err)
	EqualProduct(t, TestProduct, DBProduct)
}

func EqualProduct(t *testing.T, TestProduct models.Product, DBProduct models.Product) {
	assert.Equal(t, TestProduct.Name, DBProduct.Name)
	assert.Equal(t, TestProduct.Description, DBProduct.Description)
	assert.Equal(t, TestProduct.Parameters, DBProduct.Parameters)
	for i := range TestProduct.Images {
		assert.Equal(t, TestProduct.Images[i].Name, DBProduct.Images[i].Name)
		assert.Equal(t, TestProduct.Images[i].Key, DBProduct.Images[i].Key)

	}

}

func ReadAllProduct(t *testing.T, TestProducts []models.Product, db dbwork.DataBase) {
	DBProducts, err := db.ReadListProduct()
	assert.NoError(t, err)

	for i, DBProduct := range DBProducts {
		EqualProduct(t, TestProducts[i], DBProduct)
	}
}

func DeleteAndRead(t *testing.T, db dbwork.DataBase) {
	products, err := db.ReadListProduct()
	assert.NoError(t, err)
	count := len(products)
	id := -1
	if count >= 1 {
		id = products[0].ID
	}

	_, err = db.DeleteProduct(id)
	assert.NoError(t, err)

	products, err = db.ReadListProduct()
	assert.Equal(t, count, len(products)+1)
}

func ChangeCountProduct(t *testing.T, db dbwork.DataBase) {
	products, err := db.ReadListProduct()
	assert.NoError(t, err)

	id := 1
	count := -1
	if len(products) >= 1 {
		id = products[0].ID
		count = products[0].Count
	}

	err = db.ChangeCountProduct(id, 200)

	assert.NoError(t, err)

	product, err := db.ReadProduct(id)
	assert.NoError(t, err)

	assert.Equal(t, count+200, product.Count)

	err = db.ChangeCountProduct(id, -1-product.Count)
	assert.Error(t, err)

	count = product.Count
	err = db.ChangeCountProduct(id, -10)

	product, err = db.ReadProduct(id)
	assert.NoError(t, err)

	assert.Equal(t, product.Count, count-10)

}
