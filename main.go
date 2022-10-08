package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jvmistica/knowledge-base-go/pkg/record"
)

const apiVersion = "/api/v1"

var (
	db  *gorm.DB
	err error
)

func init() {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		log.Fatal("missing environment variable POSTGRES_HOST")
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		log.Fatal("missing environment variable POSTGRES_PORT")
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		log.Fatal("missing environment variable POSTGRES_USER")
	}

	password := os.Getenv("POSTGRES_PASS")
	if password == "" {
		log.Fatal("missing environment variable POSTGRES_PASS")
	}

	database := os.Getenv("POSTGRES_DB")
	if database == "" {
		log.Fatal("missing environment variable POSTGRES_DB")
	}

	// Connect to the database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, database, port)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the tables
	if err := db.AutoMigrate(&record.Note{}, &record.Recipe{}, &record.Script{}); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var seed = flag.Bool("seed", false, "set to true if you want to seed the database")
	flag.Parse()

	// Seed database
	if *seed {
		db.Create(&record.Notes)
		db.Create(&record.Recipes)
		db.Create(&record.Scripts)
	}

	handleRequests()
}

// handleRequests handles all the request to the APIs
func handleRequests() {
	r := record.NewRecord(db)

	http.HandleFunc(apiVersion+"/", record.GetHome)
	http.HandleFunc(apiVersion+"/notes", r.ListNotes)
	http.HandleFunc(apiVersion+"/recipes", r.ListRecipes)
	http.HandleFunc(apiVersion+"/scripts", r.ListScripts)

	http.HandleFunc(apiVersion+"/notes/new", r.CreateNote)
	http.HandleFunc(apiVersion+"/recipes/new", r.CreateRecipe)
	http.HandleFunc(apiVersion+"/scripts/new", r.CreateScript)

	http.HandleFunc(apiVersion+"/notes/delete", r.DeleteNote)
	http.HandleFunc(apiVersion+"/recipes/delete", r.DeleteRecipe)
	http.HandleFunc(apiVersion+"/scripts/delete", r.DeleteScript)

	log.Fatal(http.ListenAndServe(":10000", nil))
}
