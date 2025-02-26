package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Database struct {
	db     *sql.DB
	once   sync.Once
	config *DatabaseConfig
}

type DatabaseConfig struct {
	Type     string
	User     string
	Name     string
	Host     string
	Password string
	SslMode  string
}

var dbInstance *Database
var dbInstanceOnce sync.Once

// LoadEnv загружает переменные среды из .env файла
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func GetDatabaseInstance() *Database {
	dbInstanceOnce.Do(func() {
		LoadEnv()
		dbInstance = &Database{}
		dbInstance.config = &DatabaseConfig{
			Type:     os.Getenv("DATABASE_TYPE"),
			User:     os.Getenv("DATABASE_USER"),
			Password: os.Getenv("DATABASE_PASSWORD"),
			Name:     os.Getenv("DATABASE_NAME"),
			Host:     os.Getenv("DATABASE_HOST"),
			SslMode:  "disable",
		}
		dbInstance.connect()
	})
	return dbInstance
}

func (d *Database) connect() {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		d.config.Host, d.config.User, d.config.Password, d.config.Name, d.config.SslMode)

	var err error
	d.db, err = sql.Open("postgres", connStr)

	fmt.Println("Database connection opened succesfully")

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
}

// GetConnection returns a database connection
func (d *Database) GetConnection() *sql.DB {
	return d.db
}

// CloseAll closes the database connection pool
func (d *Database) CloseAll() {
	if err := d.db.Close(); err != nil {
		log.Fatal("Error closing the database: ", err)
	}
}

func (d *Database) ReleaseConnection() {

	if d.db != nil {
		if err := d.db.Close(); err != nil {
			log.Fatal("Error closing database connection: ", err)
		}
		fmt.Println("Database connection released successfully!")
	}
}
