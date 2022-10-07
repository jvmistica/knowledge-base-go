package record

import (
	"fmt"
	"log"
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
	var notes []Note
	if res := re.DB.Find(&notes); res.Error != nil {
		http.Error(w, fmt.Sprintf("%s", res.Error), http.StatusBadRequest)
		return
	}

	var records string
	if len(notes) > 0 {
		for _, n := range notes {
			records += fmt.Sprintf("<b>%s</b></br>%s</br></br>", n.Title, n.Content)
		}
	}

	if _, err := w.Write([]byte(records)); err != nil {
		log.Fatal(err)
	}
}

// CreateNote creates a new note
func (re *Record) CreateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Invalid method`))
	}
}
