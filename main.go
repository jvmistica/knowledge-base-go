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

// handleRequests handles all the request to the APIs
func handleRequests(r *record.Record) {
	http.HandleFunc("/", record.GetHome)

	http.HandleFunc("/notes", r.ListNotes)
	http.HandleFunc("/recipes", r.ListRecipes)
	http.HandleFunc("/scripts", r.ListScripts)

	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASS")
	database := os.Getenv("POSTGRES_DB")

	var seed = flag.Bool("seed", false, "set to true if you want to seed the database")
	flag.Parse()

	// Connect to the database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, database, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Migrate the tables
	db.AutoMigrate(&record.Note{}, &record.Recipe{}, &record.Script{})
	r := record.NewRecord(db)

	// Seed database
	if *seed {
		notes := []record.Note{
			{
				Name:        "Sample note 1",
				Description: "Some description about sample note 1",
			},
			{
				Name:        "Sample note 2",
				Description: "Some description about sample note 2",
			},
			{
				Name:        "Sample note 3",
				Description: "Some description about sample note 3",
			},
		}
		db.Create(&notes)

		recipes := []record.Recipe{
			{
				Name:        "Adobo",
				Description: "A meat dish with soy sauce, vinegar, garlic, and peppercorns.",
			},
			{
				Name:        "Rice ball",
				Description: "A simple snack made of rice, seaweed, and fillings.",
			},
			{
				Name:        "Chicken curry",
				Description: "A chicken dish with potatoes, carrots, and breaded fried chicken.",
			},
		}
		db.Create(&recipes)

		scripts := []record.Script{
			{
				Name:        "Sample script 1",
				Description: "Some description about sample script 1",
			},
			{
				Name:        "Sample script 2",
				Description: "Some description about sample script 2",
			},
			{
				Name:        "Sample script 3",
				Description: "Some description about sample script 3",
			},
		}
		db.Create(&scripts)
	}

	handleRequests(r)
}
