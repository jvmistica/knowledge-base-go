package recipe

import (
	"fmt"
	"net/http"
)

func NewRecipe(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the new recipe page for the Knowledge Base web application.")
	fmt.Println("Endpoint Hit: newRecipe")
}
