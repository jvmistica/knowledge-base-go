package main

import (
    dashboard "./dashboard"
    note "./note"
    recipe "./recipe"
    search "./search"
    script "./script"
    "fmt"
    "log"
    "net/http"
)

func handleRequests() {
    http.HandleFunc("/", homePage)
    http.HandleFunc("/dashboards", dashboardsPage)
    http.HandleFunc("/dashboards/new", dashboard.NewDashboard)
    http.HandleFunc("/notes", notesPage)
    http.HandleFunc("/notes/new", note.NewNote)
    http.HandleFunc("/recipes", recipesPage)
    http.HandleFunc("/recipes/new", recipe.NewRecipe)
    http.HandleFunc("/scripts", scriptsPage)
    http.HandleFunc("/scripts/new", script.NewScript)
    http.HandleFunc("/searches", searchesPage)
    http.HandleFunc("/searches/new", search.NewSearch)
    log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
    handleRequests()
}

// Modules
func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Welcome to the HomePage!")
    fmt.Println("Endpoint Hit: homePage")
}

func dashboardsPage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "This is the dashboards page for the Knowledge Base web application.")
    fmt.Println("Endpoint Hit: dashboardsPage")
}

func notesPage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "This is the notes page for the Knowledge Base web application.")
    fmt.Println("Endpoint Hit: notesPage")
}

func recipesPage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "This is the recipes page for the Knowledge Base web application.")
    fmt.Println("Endpoint Hit: recipesPage")
}

func scriptsPage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "This is the scripts page for the Knowledge Base web application.")
    fmt.Println("Endpoint Hit: scriptsPage")
}

func searchesPage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "This is the searches page for the Knowledge Base web application.")
    fmt.Println("Endpoint Hit: searchesPage")
}
