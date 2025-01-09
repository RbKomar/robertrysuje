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

type PageData struct {
	Title       string
	Description string
	SocialLinks []SocialLink
}

type SocialLink struct {
	Platform string
	URL      string
	Icon     string
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/gallery", handleGallery)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := PageData{
		Title:       "Portfolio Graficzne",
		Description: "Portfolio graficzne Robirysuje.pl",
		SocialLinks: []SocialLink{
			{
				Platform: "Instagram",
				URL:      "https://instagram.com/your-profile",
				Icon:     "/static/images/instagram.png",
			},
			{
				Platform: "Allegro",
				URL:      "https://allegro.pl/your-shop",
				Icon:     "/static/images/allegro.png",
			},
		},
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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
				Title: strings.Split(strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())), "_")[1],
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
