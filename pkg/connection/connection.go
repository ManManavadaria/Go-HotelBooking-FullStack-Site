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

	str := os.Getenv("CONN_STR")
	connectionStr2 := fmt.Sprintf("%s", str)

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
	createTables(db)
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

func createTables(db *sql.DB) error {
	query1 := `CREATE SEQUENCE IF NOT EXISTS reservations_id_seq;
		CREATE TABLE IF NOT EXISTS public.reservations (
		id integer NOT NULL DEFAULT nextval('reservations_id_seq'::regclass),
		first_name text COLLATE pg_catalog."default",
		last_name text COLLATE pg_catalog."default",
		email text COLLATE pg_catalog."default",
		phone text COLLATE pg_catalog."default",
		start_date timestamp without time zone,
		end_date timestamp without time zone,
		created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
		room_id integer,
		CONSTRAINT reservations_pkey PRIMARY KEY (id),
		CONSTRAINT reservations_room_id_fkey FOREIGN KEY (room_id)
			REFERENCES rooms (id) MATCH SIMPLE
			ON UPDATE NO ACTION
			ON DELETE NO ACTION
			);`

	query2 := `
			CREATE SEQUENCE IF NOT EXISTS roomtype_id_seq;
				CREATE TABLE IF NOT EXISTS public.roomtype (
					id integer NOT NULL DEFAULT nextval('roomtype_id_seq'::regclass),
					name character varying(100) COLLATE pg_catalog."default",
					CONSTRAINT roomtype_pkey PRIMARY KEY (id)
				);

			CREATE SEQUENCE IF NOT EXISTS rooms_id_seq;
			CREATE TABLE IF NOT EXISTS rooms (
				id integer NOT NULL DEFAULT nextval('rooms_id_seq'::regclass),
				roomnumber character varying(10) COLLATE pg_catalog."default",
				roomtypeid integer,
				CONSTRAINT rooms_pkey PRIMARY KEY (id),
				CONSTRAINT rooms_roomtypeid_fkey FOREIGN KEY (roomtypeid)
				REFERENCES public.roomtype (id) MATCH SIMPLE
				ON UPDATE NO ACTION
				ON DELETE NO ACTION
				);
				
				
				
	CREATE SEQUENCE IF NOT EXISTS users_id_seq;
	CREATE TABLE IF NOT EXISTS public.users (
		id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
		firstname character varying(50) COLLATE pg_catalog."default",
		lastname character varying(50) COLLATE pg_catalog."default",
		email character varying(100) COLLATE pg_catalog."default",
		password character varying(100) COLLATE pg_catalog."default",
		accesslevel integer DEFAULT 0,
		createdat timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
		updatedat timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT users_pkey PRIMARY KEY (id),
		CONSTRAINT users_email_key UNIQUE (email)
	);`

	_, err := db.Exec(query1)

	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = db.Exec(query2)
	if err != nil {
		log.Fatal(err)
		return err
	}

	var count2 int
	err = db.QueryRow("SELECT COUNT(*) FROM public.roomtype").Scan(&count2)
	if err != nil {
		log.Fatal(err)
	}

	// If the table is empty, execute the insert query
	if count2 == 0 {
		_, err := db.Exec(`
			INSERT INTO public.roomtype(name) 
			VALUES 
			('Junior Suite'),
			('Executive Suite'),
			('Super Deluxe')
		`)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Data inserted successfully.")
	} else {
		fmt.Println("Table is not empty, skipping data insertion.")
	}

	var count1 int
	err = db.QueryRow("SELECT COUNT(*) FROM rooms").Scan(&count1)
	if err != nil {
		log.Fatal(err)
	}

	// If the table is empty, execute the insert query
	if count1 == 0 {
		_, err := db.Exec(`
			INSERT INTO rooms(roomnumber, roomtypeid) 
			VALUES 
			(101,1), (102,1), (103,1), (104,1), (105,1), (106,1), 
			(201,2), (202,2), (203,2), (204,2), (205,2), (206,2), 
			(301,3), (302,3), (303,3), (304,3), (305,3), (306,3)
		`)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Data inserted successfully.")
	} else {
		fmt.Println("Table is not empty, skipping data insertion.")
	}
	return nil
}
