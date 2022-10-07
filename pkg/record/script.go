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
	for _, s := range scripts {
		records += fmt.Sprintf("%s\n", s.Name)
	}

	if _, err := w.Write([]byte(records)); err != nil {
		log.Fatal(err)
	}
}
