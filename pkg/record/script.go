package record

import (
	"fmt"
	"net/http"
	"time"
)

type Script struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (re *Record) GetScripts(w http.ResponseWriter, r *http.Request) {
	var scripts []Script
	_ = re.DB.Find(&scripts)

	var records string
	for _, s := range scripts {
		records += fmt.Sprintf("%s\n", s.Name)
	}

	w.Write([]byte(records))
}
