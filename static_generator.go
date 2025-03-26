package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
)

type StaticPageData struct {
	Title string
}

func generateStaticSite() {
	// Create the output directory structure
	setupOutputDirectories()

	// Generate home page
	generateStaticPage("index.html", "templates/layout.html", "templates/home.html", StaticPageData{Title: "model-renderer"})

	// Generate features page
	generateStaticPage("features/index.html", "templates/layout.html", "templates/features.html", StaticPageData{Title: "Features - model-renderer"})

	// Generate examples page
	generateStaticPage("examples/index.html", "templates/layout.html", "templates/examples.html", StaticPageData{Title: "Examples - model-renderer"})

	// Generate subproject pages
	generateStaticPage("subprojects/mr-graphics/index.html", "templates/layout.html", "templates/subprojects/mr-graphics.html", StaticPageData{Title: "mr-graphics - model-renderer"})
	generateStaticPage("subprojects/mr-importer/index.html", "templates/layout.html", "templates/subprojects/mr-importer.html", StaticPageData{Title: "mr-importer - model-renderer"})
	generateStaticPage("subprojects/mr-contractor/index.html", "templates/layout.html", "templates/subprojects/mr-contractor.html", StaticPageData{Title: "mr-contractor - model-renderer"})
	generateStaticPage("subprojects/mr-math/index.html", "templates/layout.html", "templates/subprojects/mr-math.html", StaticPageData{Title: "mr-math - model-renderer"})

	// Create simple redirects for docs and download
	generateStaticRedirect("docs/index.html", "/")
	generateStaticRedirect("download/index.html", "/")

	log.Println("Static website generation complete")
	log.Println("Files generated in the 'docs' directory")
	log.Println("You can now commit this to GitHub and enable GitHub Pages in your repository settings")
}

// setupOutputDirectories creates the necessary directory structure
func setupOutputDirectories() {
	dirs := []string{
		"docs",
		"docs/features",
		"docs/examples",
		"docs/docs",
		"docs/download",
		"docs/subprojects/mr-graphics",
		"docs/subprojects/mr-importer",
		"docs/subprojects/mr-contractor",
		"docs/subprojects/mr-math",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Copy assets directory
	if err := copyAssetDir("assets", "docs/assets"); err != nil {
		log.Fatalf("Failed to copy assets: %v", err)
	}
}

// generateStaticPage renders a template to a static HTML file
func generateStaticPage(outputPath, layoutPath, contentPath string, data StaticPageData) {
	// Parse templates
	tmpl, err := template.ParseFiles(layoutPath, contentPath)
	if err != nil {
		log.Fatalf("Failed to parse templates %s and %s: %v", layoutPath, contentPath, err)
	}

	// Create output file
	outputFile := filepath.Join("docs", outputPath)
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", outputFile, err)
	}
	defer file.Close()

	// Execute template and write to file
	if err := tmpl.Execute(file, data); err != nil {
		log.Fatalf("Failed to execute template for %s: %v", outputPath, err)
	}

	log.Printf("Generated %s", outputFile)
}

// generateStaticRedirect creates a simple HTML redirect page
func generateStaticRedirect(outputPath, target string) {
	// Create redirect HTML
	redirectHTML := `<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="refresh" content="0; url=` + target + `">
</head>
<body>
    <p>Redirecting to <a href="` + target + `">` + target + `</a></p>
</body>
</html>`

	// Create output file
	outputFile := filepath.Join("docs", outputPath)
	err := os.WriteFile(outputFile, []byte(redirectHTML), 0644)
	if err != nil {
		log.Fatalf("Failed to create redirect file %s: %v", outputFile, err)
	}

	log.Printf("Generated redirect from %s to %s", outputFile, target)
}

// copyAssetDir recursively copies a directory tree
func copyAssetDir(src string, dst string) error {
	// Create destination directory
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	// Get directory contents
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories
			if err := copyAssetDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}

			if err := os.WriteFile(dstPath, data, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}
