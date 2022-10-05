package record

import (
	"fmt"
	"net/http"
)

// GetHome returns the contents of the homepage
func GetHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
