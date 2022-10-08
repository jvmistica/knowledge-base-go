package record

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Recipe is the structure of the recipes table
type Recipe struct {
	ID          uint
	Name        string
	Description string
	Instruction string
	Category    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ListRecipes lists all the recipes in the database
func (re *Record) ListRecipes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var recipes []Recipe
	if res := re.DB.Find(&recipes); res.Error != nil {
		http.Error(w, fmt.Sprintf("%s", res.Error), http.StatusInternalServerError)
		return
	}

	var records string
	if len(recipes) > 0 {
		for _, r := range recipes {
			records += fmt.Sprintf("<b>%s</b></br>%s</br></br>", r.Name, r.Description)
		}
	}

	if _, err := w.Write([]byte(records)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateRecipe creates a new recipe
func (re *Record) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var recipe Recipe
	if err := json.Unmarshal(body, &recipe); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := re.DB.Create(&recipe)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("%s", result.Error), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
