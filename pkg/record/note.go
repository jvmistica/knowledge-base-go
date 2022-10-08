package record

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Note is the structure of the notes table
type Note struct {
	ID        uint
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ListNotes lists all the notes in the database
func (re *Record) ListNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var notes []Note
	if res := re.DB.Find(&notes); res.Error != nil {
		http.Error(w, fmt.Sprintf("%s", res.Error), http.StatusInternalServerError)
		return
	}

	var records string
	if len(notes) > 0 {
		for _, n := range notes {
			records += fmt.Sprintf("<b>%s</b></br>%s</br></br>", n.Title, n.Content)
		}
	}

	if _, err := w.Write([]byte(records)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// CreateNote creates a new note
func (re *Record) CreateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var note Note
	if err := json.Unmarshal(body, &note); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := re.DB.Create(&note)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("%s", result.Error), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
