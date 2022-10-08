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
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Instruction string    `json:"instruction"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

	recipesList, err := json.Marshal(recipes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(recipesList)
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
	defer r.Body.Close()

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

// DeleteRecipe deletes a recipe
func (re *Record) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing query parameter: 'id'", http.StatusBadRequest)
		return
	}

	result := re.DB.Where("id = ?", id).Delete(Recipe{})
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("%s", result.Error), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetRecipe gets the details of a specific recipe
func (re *Record) GetRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing query parameter: 'id'", http.StatusBadRequest)
		return
	}

	var recipe Recipe
	result := re.DB.Where("id = ?", id).Find(&recipe)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("%s", result.Error), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	details, err := json.Marshal(recipe)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(details)
}
