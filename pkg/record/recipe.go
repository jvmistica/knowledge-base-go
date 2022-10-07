package record

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Recipe is the structure of the recipes table
type Recipe struct {
	ID          uint
	Name        string
	Description string
	Instruction string
	Category    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ListRecipes lists all the recipes in the database
func (re *Record) ListRecipes(w http.ResponseWriter, r *http.Request) {
	var recipes []Recipe
	if res := re.DB.Find(&recipes); res.Error != nil {
		http.Error(w, fmt.Sprintf("%s", res.Error), http.StatusBadRequest)
		return
	}

	var records string
	if len(recipes) > 0 {
		for _, r := range recipes {
			records += fmt.Sprintf("<b>%s</b></br>%s</br></br>", r.Name, r.Description)
		}
	}

	if _, err := w.Write([]byte(records)); err != nil {
		log.Fatal(err)
	}
}
