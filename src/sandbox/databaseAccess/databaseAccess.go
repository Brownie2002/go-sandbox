package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

// stackTracer helps to display the stack trace when errors happen.
type stackTracer interface {
	StackTrace() errors.StackTrace
}

// showStackTrace is there to show ho to display a stacktrace in Golang
// using the library pkg.Errors from Dave Cheney.
func printStackTrace(err error) {
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			fmt.Printf("%+s: %d\n", f, f)
		}
	}
	fmt.Printf("%v\n", err)
}

// initDb initializes a SQLlite3 database with a file set on disk. The
// path to the file is defined as the parameter of the function.
// Warning: this file is deleted if it exists.
// The function returns a pointer to an opened connection to the database.
// If opening the database fails, it returns an error.
func initDb(dbfile string) (*sql.DB, error) {

	// Remove the existing file if exists
	err := os.Remove(dbfile)
	if err != nil {
		log.Printf("There is no existing database '%s'. Nothing to delete.\n", dbfile)
	}

	// Init the db
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open the database.")
	}

	return db, nil
}

// createSampleTable simply creates a predefined table name foo with 2 columns:
// an id, as an integer, and name as text. Syntax of creation should meet the SQL
// syntax to be runnable on multiple databases.
func createSampleTable(db *sql.DB) error {

	sqlStmt := `create table foo (id integer not null primary key, name text);`

	_, err := db.Exec(sqlStmt)
	if err != nil {
		return errors.Wrap(err, "Failed to create the sample table in the database.")
	}

	return nil
}

// deleteSampleTableIfExists is the pending method to createSampleTable. It
// drops the foo table if it exists. This allows to run the code multiple times
// on the same database with no error when attempting to create the table foo.
func deleteSampleTableIfExists(db *sql.DB) error {
	fmt.Println("Dropping table foo if it exists...")
	stmt := "drop table IF EXISTS foo"
	_, err := db.Exec(stmt)
	if err != nil {
		return errors.Wrap(err, "Failed to execute drop statement")
	}
	return nil
}

// insertSampleData insert 10 rows in the table foo. It demonstates the use
// of prepared statement. Note that we use the syntax $i for the variable
// parameters. This may not work for all databases. This syntax works for
// both SQLlite 3 and PostgreSQL.
func insertSampleData(db *sql.DB) error {
	fmt.Println("Inserting sample data (10 rows)")
	tx, err := db.Begin()
	if err != nil {
		return errors.Wrap(err, "Failed to create a transaction..")
	}
	stmt, err := tx.Prepare("insert into foo(id, name) values ($1,$2);")
	if err != nil {
		return errors.Wrap(err, "Failed to create prepared statement.")
	}
	defer stmt.Close()

	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("nom%03d", i)
		_, err = stmt.Exec(i, fmt.Sprintf("%03d", i, name))

		if err != nil {
			return errors.Wrap(err, "Failed to insert sample data in the sample table.")
		}
	}
	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "Failed to commit the transaction.")
	}

	return nil
}

// selectDataInSampleTable selects all the rows in the foo table.
// Then it displays them on stdout.
func selectDataInSampleTable(db *sql.DB) error {
	rows, err := db.Query("select id, name from foo;")
	if err != nil {
		return errors.Wrap(err, "Failed to prepare select query.")
	}
	defer rows.Close()

	for rows.Next() {

		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			return errors.Wrap(err, "failed to scan output of query.")
		}
		fmt.Println(id, name)
	}

	err = rows.Err()
	if err != nil {
		return errors.Wrap(err, "error in the select query.")
	}

	return nil
}

// cleanupSampleTabledb deletes all the rows from the foo table that is
// driving this test.
func cleanupSampleTabledb(db *sql.DB) error {
	_, err := db.Exec("delete from foo;")
	if err != nil {
		return errors.Wrap(err, "failed to delete rows from the table foo")
	}

	return nil
}

// insertDataAndSelectThem demontrates how to insert multiple rows with the
// same query. Then the inserted data is selected with a SQL query and
// displayed on the stdout.
func insertDataAndSelectThem(db *sql.DB) error {
	_, err := db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz');")
	if err != nil {
		return errors.Wrap(err, "failed to insert data into table foo")
	}

	rows, err := db.Query("select id, name from foo;")
	if err != nil {
		return errors.Wrap(err, "failed to query table foo")
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			return errors.Wrap(err, "Failed to read values from the row.")
		}
		fmt.Println(id, name)
	}
	err = rows.Err()
	if err != nil {
		return errors.Wrap(err, "Error happened while reading the rows.")
	}

	return nil
}

// connectToPostgres uses a connection string that details how to access to
// a postgres database.
func connectToPostgres(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to connect to the database with parameters '"+connStr+"'")
	}
	return db, nil
}

// showExamplesWithSqlLite runs code on a SQL lite database to learn how
// it works to interact between Golang and a SQL Lite database.
func showExamplesWithSqlLite() {
	fmt.Println("Showing output from statements ran against SqlLite.")
	fmt.Println("===================================================")
	dbfile := "./foo.db"
	db, err := initDb(dbfile)
	if err != nil {
		printStackTrace(err)
		return
	}
	defer db.Close()

	err = createSampleTable(db)
	if err != nil {
		printStackTrace(err)
		return
	}

	err = insertSampleData(db)
	if err != nil {
		printStackTrace(err)
		return
	}

	err = selectDataInSampleTable(db)
	if err != nil {
		printStackTrace(err)
		return
	}

	err = cleanupSampleTabledb(db)
	if err != nil {
		printStackTrace(err)
		return
	}

	err = insertDataAndSelectThem(db)
	if err != nil {
		printStackTrace(err)
		return
	}
}

type PostgresDb struct {
	host     string
	port     int
	dbname   string
	user     string
	password string
	sslmode  string
}

// buildConnectionString helps to rebuild a connection string for a postgresql
// database that is not that obvious to build. There are probably better solution
// in the libpq (and/or more elegant).
func (p PostgresDb) buildConnectionString() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s", p.user, p.password, p.dbname, p.host, p.port, p.sslmode)
}

// showExamplesWithPostgres runs some code to use the PostgreSQL database
// from Golang code. It contains some actions (create a table, insert rows,
// delete a table, make queries on tables, use prepared statement) that are
// mainly part of the usual learning curve when switching from a language to
// another.
func showExamplesWithPostgres(dbParams PostgresDb) {
	fmt.Println("Showing output from statements ran against PostgreSQL.")
	fmt.Println("======================================================")
	connStr := dbParams.buildConnectionString()

	db, err := connectToPostgres(connStr)
	if err != nil {
		printStackTrace(err)
		return
	}
	fmt.Println("Connected to the PostgreSQL database. Proceeding")
	defer db.Close()

	err = deleteSampleTableIfExists(db)
	if err != nil {
		printStackTrace(err)
		return
	}

	err = createSampleTable(db)
	if err != nil {
		printStackTrace(err)
		return
	}

	err = insertSampleData(db)
	if err != nil {
		fmt.Println(err)
		printStackTrace(err)
		return
	}

	err = selectDataInSampleTable(db)
	if err != nil {
		printStackTrace(err)
		return
	}

	rows, err := db.Query("SELECT srid, proj4text FROM spatial_ref_sys;")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var srid int
		var proj4text string

		err := rows.Scan(&srid, &proj4text)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Test: ", srid, proj4text)

	}
}

// loadPgParams retrieves the configuration from a file on disk. The
// file is supposed to have one key per line (user, password, dbname,
// host, port, sslmode).
// The function returns a structure with those values.
func loadPgParams(pgfile string) (PostgresDb, error) {
	fmt.Printf("Using file '%s' as parameters to connect to PostgreSQL.\n", pgfile)

	f, err := os.Open(pgfile)
	out := PostgresDb{}
	if err != nil {
		return out, errors.Wrap(err, "cannot open configuration file containing informations to connect to the database")
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	var key, value string

	for s.Scan() {
		line := s.Text()
		elements := strings.Split(line, "=")
		key = elements[0]
		value = elements[1]
		switch key {
		case "user":
			out.user = value
		case "password":
			out.password = value
		case "dbname":
			out.dbname = value
		case "host":
			out.host = value
		case "port":
			out.port, err = strconv.Atoi(value)
			if err != nil {
				msg := "The port you have set '" + value + "' for the database is not an integer. Please fix the configuration."
				return out, errors.New(msg)
			}
		case "sslmode":
			out.sslmode = value
		}
	}
	return out, nil
}

// main runs examples of code for SQLlite3 and PostgreSQL databases.
func main() {

	showExamplesWithSqlLite()

	// Default configuration for PostgreSQL. This configuration matches
	// the configuration required for the container
	// docker.io/helmi03/docker-postgis
	params := PostgresDb{
		user:     "docker",
		password: "docker",
		dbname:   "postgres",
		host:     "localhost",
		port:     5432,
		sslmode:  "require",
	}

	// err must be declared here. If you use := syntax in if statement, the
	// variables do not live outside of the if branch.
	var err error
	if len(os.Args) > 1 {
		params, err = loadPgParams(os.Args[1])
		if err != nil {
			fmt.Printf("Failed to load database properties to connect: %v\n", err)
			return
		}
	}
	showExamplesWithPostgres(params)
}
