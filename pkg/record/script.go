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
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
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

	var records string
	if len(scripts) > 0 {
		for _, s := range scripts {
			records += fmt.Sprintf("<b>%s</b></br>%s</br></br>", s.Name, s.Description)
		}
	}

	if _, err := w.Write([]byte(records)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
