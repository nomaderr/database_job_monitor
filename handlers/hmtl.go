package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	filePath := "templates/index.html"
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		log.Println("Error loading HTML file:", err)
		http.Error(w, "Error loading HTML file", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Error rendering HTML file:", err)
		http.Error(w, "Error rendering HTML file", http.StatusInternalServerError)
	}
}
