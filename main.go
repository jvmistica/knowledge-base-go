package main

import (
	// "encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jvmistica/knowledge-base-go/pkg/note"
	"github.com/jvmistica/knowledge-base-go/pkg/recipe"
	"github.com/jvmistica/knowledge-base-go/pkg/script"
)

// TODO: Use struct once all controllers are implemented
var (
	DB *gorm.DB
)

type Note struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Recipe struct {
	ID          uint
	Name        string
	Description string
	Instruction string
	Category    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Script struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func handleRequests() {
	http.HandleFunc("/", homePage)

	http.HandleFunc("/notes", notesPage)
	http.HandleFunc("/notes/new", note.NewNote)

	http.HandleFunc("/recipes", recipesPage)
	http.HandleFunc("/recipes/new", recipe.NewRecipe)

	http.HandleFunc("/scripts", scriptsPage)
	http.HandleFunc("/scripts/new", script.NewScript)

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
	db.AutoMigrate(&Recipe{})

	DB = db

	handleRequests()
}

// Modules
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func notesPage(w http.ResponseWriter, r *http.Request) {
	var notes []Note
	_ = DB.Find(&notes)

	var records string
	for _, n := range notes {
		records += fmt.Sprintf("%s\n", n.Name)
	}

	w.Write([]byte(records))
}

func recipesPage(w http.ResponseWriter, r *http.Request) {
	var recipes []Recipe
	_ = DB.Find(&recipes)

	var records string
	for _, r := range recipes {
		records += fmt.Sprintf("%s\n", r.Name)
	}

	w.Write([]byte(records))
}

func scriptsPage(w http.ResponseWriter, r *http.Request) {
	var scripts []Script
	_ = DB.Find(&scripts)

	var records string
	for _, s := range scripts {
		records += fmt.Sprintf("%s\n", s.Name)
	}

	w.Write([]byte(records))
}
