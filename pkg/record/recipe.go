package record

import (
	"fmt"
	"net/http"
	"time"
)

type Recipe struct {
	ID          uint
	Name        string
	Description string
	Instruction string
	Category    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (re *Record) GetRecipes(w http.ResponseWriter, r *http.Request) {
	var recipes []Recipe
	if res := re.DB.Find(&recipes); res.Error != nil {
		http.Error(w, fmt.Sprintf("%s", res.Error), http.StatusBadRequest)
		return
	}

	var records string
	for _, r := range recipes {
		records += fmt.Sprintf("%s\n", r.Name)
	}

	w.Write([]byte(records))
}
