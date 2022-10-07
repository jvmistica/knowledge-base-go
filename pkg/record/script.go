package record

import (
	"fmt"
	"log"
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
	var scripts []Script
	if res := re.DB.Find(&scripts); res.Error != nil {
		http.Error(w, fmt.Sprintf("%s", res.Error), http.StatusBadRequest)
		return
	}

	var records string
	if len(scripts) > 0 {
		for _, s := range scripts {
			records += fmt.Sprintf("<b>%s</b></br>%s</br></br>", s.Name, s.Description)
		}
	}

	if _, err := w.Write([]byte(records)); err != nil {
		log.Fatal(err)
	}
}

// CreateScript creates a new recipe
func (re *Record) CreateScript(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Invalid method`))
	}
}
