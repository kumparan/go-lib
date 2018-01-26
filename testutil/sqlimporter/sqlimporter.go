package sqlimporter

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func createDSN() string {
	return ""
}

// Connect to database
func Connect(driver, dsn string) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	db, err := sqlx.ConnectContext(ctx, driver, dsn)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

// DBNameDefault default database name for sqlimporter
const DBNameDefault = "SQL_IMPORTER_DB_"

// CreateRandomDB for creating a database name/schema for test
func CreateRandomDB(driver, dsn string) (*sqlx.DB, func() error, error) {
	// create a new database
	// database name is always a random name
	rand.Seed(time.Now().UnixNano())
	dbName := DBNameDefault + strconv.FormatInt(rand.Int63(), 10)
	return CreateDB(driver, dbName, dsn)
}

// CreateDB used to create database
// and import all queries located in a directories
func CreateDB(driver, dbName, dsn string) (*sqlx.DB, func() error, error) {
	defaultDrop := func() error {
		return nil
	}
	db, err := Connect(driver, dsn)
	if err != nil {
		return nil, defaultDrop, err
	}

	createDBQuery := fmt.Sprintf(getDialect(driver, "create"), dbName)
	// exec create new b
	_, err = db.Exec(createDBQuery)
	if err != nil {
		return nil, defaultDrop, err
	}

	// use new db
	err = selectDB(db, dbName)
	if err != nil {
		return nil, defaultDrop, err
	}
	return db, func() error {
		deleteDatabaseQuery := fmt.Sprintf(getDialect(driver, "drop"), dbName)
		_, err := db.Exec(deleteDatabaseQuery)
		if err != nil {
			return err
		}
		return db.Close()
	}, nil
}

// selectDB for selecting database/schema
func selectDB(db *sqlx.DB, dbName string) error {
	// use new db
	useDatabaseQuery := fmt.Sprintf(getDialect(db.DriverName(), "use"), dbName)
	_, err := db.Exec(useDatabaseQuery)
	return err
}

// ImportSchemaFromFiles to import all *.sql file from directory
func ImportSchemaFromFiles(db *sqlx.DB, filepath string) error {
	files, err := getFileList(filepath)
	if err != nil {
		return err
	}

	// an sql file will be executed as one batch of transaction
	for _, file := range files {
		sqlContents, err := parseFiles(file)
		if err != nil {
			return err
		}
		// end if empty
		if len(sqlContents) == 0 {
			return nil
		}

		tx, err := db.BeginTx(context.TODO(), nil)
		if err != nil {
			return err
		}

		var query string
		for key := range sqlContents {
			query = sqlContents[key]
			_, err = tx.ExecContext(context.TODO(), query)
			if err != nil {
				break
			}
		}

		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return fmt.Errorf("Failed to rollback from file %s with error %s", file, errRollback.Error())
			}
			return fmt.Errorf("Failed to execute file %s with error %s and query: \n %s", file, err.Error(), query)
		}
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("Failed to commit from file %s with error %s", file, err.Error())
		}

	}
	return nil
}
