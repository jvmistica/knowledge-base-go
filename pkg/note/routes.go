package note

import (
	"fmt"
	"net/http"
)

func NewNote(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is the new note page for the Knowledge Base web application.")
	fmt.Println("Endpoint Hit: newNote")
}
