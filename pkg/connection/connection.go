package connection

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSQL creates database pool for Postgres
func ConnectSQL() (*DB, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	a := os.Getenv("HOST")
	b := os.Getenv("DBPORT")
	c := os.Getenv("DB_USER")
	d := os.Getenv("PASSWORD")
	e := os.Getenv("DB_NAME")
	f := os.Getenv("SSL_MODE")

	fmt.Println(a, b, c, d, e, f)

	connectionStr2 := fmt.Sprintf("postgresql://manm49061:%s@lingering-hat-17174145.ap-southeast-1.aws.neon.tech/%s?sslmode=%s", d, e, f)
	// connectionStr2 := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", a, b, c, d, e, f)

	db, err := NewDatabase(connectionStr2)
	if err != nil {
		log.Fatal(err, "db connection err")
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifetime)

	dbConn.SQL = db

	err = testDB(db)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// NewDatabase creates a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// testDB tries to ping the database
func testDB(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}
