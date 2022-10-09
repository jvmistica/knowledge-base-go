package record

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Script is the structure of the scripts table
type Script struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ListScripts lists all the scripts in the database
func (re *Record) ListScripts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var scripts []Script
	if res := re.DB.Find(&scripts); res.Error != nil {
		http.Error(w, fmt.Sprintf("%s", res.Error), http.StatusInternalServerError)
		return
	}

	scriptsList, err := json.Marshal(scripts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	w.Write(scriptsList)
}

// CreateScript creates a new recipe
func (re *Record) CreateScript(w http.ResponseWriter, r *http.Request) {
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

	var script Script
	if err := json.Unmarshal(body, &script); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := re.DB.Create(&script)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("%s", result.Error), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeleteScript deletes a script
func (re *Record) DeleteScript(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing query parameter: 'id'", http.StatusBadRequest)
		return
	}

	result := re.DB.Where("id = ?", id).Delete(Script{})
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

// GetScript gets the details of a specific script
func (re *Record) GetScript(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing query parameter: 'id'", http.StatusBadRequest)
		return
	}

	var script Script
	result := re.DB.Where("id = ?", id).Find(&script)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("%s", result.Error), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	details, err := json.Marshal(script)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(details)
}

// UpdateScript updates an existing script
func (re *Record) UpdateScript(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var script Script
	if err := json.Unmarshal(body, &script); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := re.DB.Model(&Script{}).Where("id = ?", script.ID).Updates(script)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("%s", result.Error), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
