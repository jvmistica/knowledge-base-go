package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jvmistica/knowledge-base-go/pkg/record"
)

func handleRequests(r *record.Record) {
	http.HandleFunc("/", homePage)

	http.HandleFunc("/notes", r.GetNotes)
	http.HandleFunc("/recipes", r.GetRecipes)
	http.HandleFunc("/scripts", r.GetScripts)

	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASS")
	database := os.Getenv("POSTGRES_DB")

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

	handleRequests(r)
}

// Modules
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
