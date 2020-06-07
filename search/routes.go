package search

import (
    "fmt"
    "net/http"
)


func NewSearch(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "This is the new search page for the Knowledge Base web application.")
    fmt.Println("Endpoint Hit: newSearch")
}
