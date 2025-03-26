package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type PageData struct {
	Title string
}

func main() {
	// Parse command line flags
	githubPages := flag.Bool("github-pages", false, "Generate GitHub Pages static site")
	port := flag.String("port", "4747", "Port to run the server on")
	flag.Parse()

	// If --github-pages flag is set, generate GitHub Pages site and exit
	if *githubPages {
		GenerateGitHubPages()
		return
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Serve static files
	fileServer := http.FileServer(http.Dir("assets"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fileServer))

	// Routes
	r.Get("/", handleHome)
	r.Get("/features", handleFeatures)
	r.Get("/examples", handleExamples)
	r.Get("/docs", handleDocs)
	r.Get("/download", handleDownload)

	// Subproject routes
	r.Get("/subprojects/mr-graphics", handleMRGraphics)
	r.Get("/subprojects/mr-importer", handleMRImporter)
	r.Get("/subprojects/mr-contractor", handleMRContractor)
	r.Get("/subprojects/mr-math", handleMRMath)

	// Start server
	log.Printf("Server starting on :%s - Visit http://localhost:%s\n", *port, *port)
	log.Printf("To generate GitHub Pages site, restart with: go run . --github-pages\n")
	if err := http.ListenAndServe(":"+*port, r); err != nil {
		log.Fatal(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
	data := PageData{
		Title: "Game Engine",
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleFeatures(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/features.html"))
	data := PageData{
		Title: "Features - Game Engine",
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleExamples(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/examples.html"))
	data := PageData{
		Title: "Examples - Game Engine",
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleDocs(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // Temporary redirect to home
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect) // Temporary redirect to home
}

func handleMRGraphics(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/subprojects/mr-graphics.html"))
	data := PageData{
		Title: "MR Graphics - Game Engine",
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleMRImporter(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/subprojects/mr-importer.html"))
	data := PageData{
		Title: "MR Importer - Game Engine",
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleMRContractor(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/subprojects/mr-contractor.html"))
	data := PageData{
		Title: "MR Contractor - Game Engine",
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleMRMath(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/subprojects/mr-math.html"))
	data := PageData{
		Title: "MR Math - Game Engine",
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
