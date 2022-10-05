package record

import (
	"fmt"
	"net/http"
	"time"
)

// Note is the structure of the notes table
type Note struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ListNotes lists all the notes in the database
func (re *Record) ListNotes(w http.ResponseWriter, r *http.Request) {
	var notes []Note
	if res := re.DB.Find(&notes); res.Error != nil {
		http.Error(w, fmt.Sprintf("%s", res.Error), http.StatusBadRequest)
		return
	}

	var records string
	for _, n := range notes {
		records += fmt.Sprintf("%s\n", n.Name)
	}

	w.Write([]byte(records))
}
