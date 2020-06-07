package script

import (
    "fmt"
    "net/http"
)


func NewScript(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "This is the new script page for the Knowledge Base web application.")
    fmt.Println("Endpoint Hit: newScript")
}
