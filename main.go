package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type PageData struct {
	Title   string
	Lang    string
	Year    int
	BaseURL string
}

// Load all templates at startup instead of on each request
var templates map[string]*template.Template

func loadTemplates() {
	templates = make(map[string]*template.Template)

	// Define base templates that should be included in every page
	baseTemplates := []string{"templates/layout.html", "templates/translations.html"}

	// Load page templates
	templates["home"] = template.Must(template.ParseFiles(append(baseTemplates, "templates/home.html")...))
	templates["features"] = template.Must(template.ParseFiles(append(baseTemplates, "templates/features.html")...))
	templates["examples"] = template.Must(template.ParseFiles(append(baseTemplates, "templates/examples.html")...))
	templates["mr-graphics"] = template.Must(template.ParseFiles(append(baseTemplates, "templates/subprojects/mr-graphics.html")...))
	templates["mr-importer"] = template.Must(template.ParseFiles(append(baseTemplates, "templates/subprojects/mr-importer.html")...))
	templates["mr-contractor"] = template.Must(template.ParseFiles(append(baseTemplates, "templates/subprojects/mr-contractor.html")...))
	templates["mr-math"] = template.Must(template.ParseFiles(append(baseTemplates, "templates/subprojects/mr-math.html")...))
}

func main() {
	// Parse command line flags
	githubPages := flag.Bool("github-pages", false, "Generate GitHub Pages static site")
	port := flag.String("port", "4747", "Port to run the server on")
	flag.Parse()

	// Load all templates
	loadTemplates()

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

func getLanguage(r *http.Request) string {
	lang := r.URL.Query().Get("lang")
	if lang == "ru" {
		return "ru"
	}
	return "en" // Default to English
}

func getPageData(title string, r *http.Request) PageData {
	return PageData{
		Title: title,
		Lang:  getLanguage(r),
		Year:  time.Now().Year(),
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	if t, ok := templates[tmpl]; ok {
		err := t.ExecuteTemplate(w, "layout", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error rendering template: %v", err)
		}
	} else {
		http.Error(w, "Template not found", http.StatusInternalServerError)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	data := getPageData("Game Engine", r)
	renderTemplate(w, "home", data)
}

func handleFeatures(w http.ResponseWriter, r *http.Request) {
	data := getPageData("Features - Game Engine", r)
	renderTemplate(w, "features", data)
}

func handleExamples(w http.ResponseWriter, r *http.Request) {
	data := getPageData("Examples - Game Engine", r)
	renderTemplate(w, "examples", data)
}

func handleDocs(w http.ResponseWriter, r *http.Request) {
	// Preserve language parameter when redirecting to home
	redirectURL := "/"
	if lang := r.URL.Query().Get("lang"); lang == "ru" {
		redirectURL = "/?lang=ru"
	}
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	// Preserve language parameter when redirecting to home
	redirectURL := "/"
	if lang := r.URL.Query().Get("lang"); lang == "ru" {
		redirectURL = "/?lang=ru"
	}
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

func handleMRGraphics(w http.ResponseWriter, r *http.Request) {
	data := getPageData("MR Graphics - Game Engine", r)
	renderTemplate(w, "mr-graphics", data)
}

func handleMRImporter(w http.ResponseWriter, r *http.Request) {
	data := getPageData("MR Importer - Game Engine", r)
	renderTemplate(w, "mr-importer", data)
}

func handleMRContractor(w http.ResponseWriter, r *http.Request) {
	data := getPageData("MR Contractor - Game Engine", r)
	renderTemplate(w, "mr-contractor", data)
}

func handleMRMath(w http.ResponseWriter, r *http.Request) {
	data := getPageData("MR Math - Game Engine", r)
	renderTemplate(w, "mr-math", data)
}
