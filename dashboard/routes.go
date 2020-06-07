package dashboard

import (
    "fmt"
    "net/http"
)


func NewDashboard(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "This is the new dashboard page for the Knowledge Base web application.")
    fmt.Println("Endpoint Hit: newDashboard")
}
