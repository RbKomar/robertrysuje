package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	Path  string
	Title string
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/gallery", handleGallery)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func handleGallery(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("static/gallery")
	if err != nil {
		log.Printf("Error reading gallery: %v", err)
		http.Error(w, "Error loading gallery", http.StatusInternalServerError)
		return
	}

	var images []Image
	for _, file := range files {
		if !file.IsDir() && isImageFile(file.Name()) {
			images = append(images, Image{
				Path:  "/static/gallery/" + file.Name(),
				Title: strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())),
			})
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/gallery.html"))
	tmpl.Execute(w, images)
}

func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}
