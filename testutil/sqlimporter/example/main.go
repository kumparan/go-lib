package main

// ONLY EXAMPLE, COMMENTED BECAUSE OF MISSING DEPENDENCY github.com/lab46/example/pkg/testutil/sqlimporter
//go:binary-only-package

// import (
// 	"log"

// 	// "github.com/lab46/example/pkg/testutil/sqlimporter"
// 	_ "github.com/lib/pq"
// )

func main() {
	// 	dsn := "postgres://exampleapp:exampleapp@localhost:5432?sslmode=disable"
	// 	// dsn := "user=exampleapp password=exampleapp host=127.0.0.1:5432 dbname=exampleapp sslmode=disable"
	// 	// db, err := sqlx.Open("postgres", dsn)
	// 	db, drop, err := sqlimporter.CreateRandomDB("postgres", dsn)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	err = sqlimporter.ImportSchemaFromFiles(db, "../files")
	// 	if err != nil {
	// 		log.Fatal("Failed to import ", err.Error())
	// 	}
	// 	defer func() {
	// 		err := drop()
	// 		if err != nil {
	// 			log.Printf("Failed to drop. Error: %s", err.Error())
	// 		}
	// 	}()
	// 	err = db.Ping()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Println("finished")
}
