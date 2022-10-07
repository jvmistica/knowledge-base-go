package record

import (
	"fmt"
	"log"
	"net/http"
)

// GetHome returns the contents of the homepage
func GetHome(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "Welcome to the HomePage!"); err != nil {
		log.Fatal(err)
	}
}
