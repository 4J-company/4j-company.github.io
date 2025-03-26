package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
)

// GitHubPagesGenerator handles building static files for GitHub Pages
type GitHubPagesGenerator struct {
	OutputDir string
}

// NewGitHubPagesGenerator creates a new generator instance
func NewGitHubPagesGenerator() *GitHubPagesGenerator {
	return &GitHubPagesGenerator{
		OutputDir: "docs",
	}
}

// Run executes the static site generation process
func (g *GitHubPagesGenerator) Run() {
	fmt.Println("Building static site for GitHub Pages...")

	g.setupDirectories()
	g.generateAllPages()

	fmt.Println("\nStatic site generation complete!")
	fmt.Println("\nTo deploy to GitHub Pages:")
	fmt.Println("1. Create a GitHub repository")
	fmt.Println("2. Commit and push your code including the 'docs' directory")
	fmt.Println("3. Go to repository Settings -> Pages")
	fmt.Println("4. Under 'Source', select 'Deploy from a branch'")
	fmt.Println("5. Select 'main' branch and '/docs' folder, then click 'Save'")
	fmt.Println("\nYour site will be available at https://yourusername.github.io/repository-name/")
}

// setupDirectories creates the necessary directory structure
func (g *GitHubPagesGenerator) setupDirectories() {
	dirs := []string{
		g.OutputDir,
		filepath.Join(g.OutputDir, "features"),
		filepath.Join(g.OutputDir, "examples"),
		filepath.Join(g.OutputDir, "docs"),
		filepath.Join(g.OutputDir, "download"),
		filepath.Join(g.OutputDir, "subprojects/mr-graphics"),
		filepath.Join(g.OutputDir, "subprojects/mr-importer"),
		filepath.Join(g.OutputDir, "subprojects/mr-contractor"),
		filepath.Join(g.OutputDir, "subprojects/mr-math"),
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create empty .nojekyll file to disable Jekyll processing
	noJekyllPath := filepath.Join(g.OutputDir, ".nojekyll")
	if err := os.WriteFile(noJekyllPath, []byte{}, 0644); err != nil {
		log.Fatalf("Failed to create .nojekyll file: %v", err)
	}

	// Copy assets directory
	if err := g.copyDirectory("assets", filepath.Join(g.OutputDir, "assets")); err != nil {
		log.Fatalf("Failed to copy assets: %v", err)
	}
}

// generateAllPages generates all static HTML pages
func (g *GitHubPagesGenerator) generateAllPages() {
	// Generate home page
	g.generatePage("index.html", "templates/layout.html", "templates/home.html", PageData{Title: "model-renderer"})

	// Generate features page
	g.generatePage("features/index.html", "templates/layout.html", "templates/features.html", PageData{Title: "Features - model-renderer"})

	// Generate examples page
	g.generatePage("examples/index.html", "templates/layout.html", "templates/examples.html", PageData{Title: "Examples - model-renderer"})

	// Generate subproject pages
	g.generatePage("subprojects/mr-graphics/index.html", "templates/layout.html", "templates/subprojects/mr-graphics.html", PageData{Title: "mr-graphics - model-renderer"})
	g.generatePage("subprojects/mr-importer/index.html", "templates/layout.html", "templates/subprojects/mr-importer.html", PageData{Title: "mr-importer - model-renderer"})
	g.generatePage("subprojects/mr-contractor/index.html", "templates/layout.html", "templates/subprojects/mr-contractor.html", PageData{Title: "mr-contractor - model-renderer"})
	g.generatePage("subprojects/mr-math/index.html", "templates/layout.html", "templates/subprojects/mr-math.html", PageData{Title: "mr-math - model-renderer"})

	// Create simple redirects for docs and download
	g.generateRedirect("docs/index.html", "/")
	g.generateRedirect("download/index.html", "/")
}

// generatePage renders a template to a static HTML file
func (g *GitHubPagesGenerator) generatePage(outputPath, layoutPath, contentPath string, data PageData) {
	// Parse templates
	tmpl, err := template.ParseFiles(layoutPath, contentPath)
	if err != nil {
		log.Fatalf("Failed to parse templates %s and %s: %v", layoutPath, contentPath, err)
	}

	// Create output file
	outputFile := filepath.Join(g.OutputDir, outputPath)
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", outputFile, err)
	}
	defer file.Close()

	// Execute template and write to file
	if err := tmpl.Execute(file, data); err != nil {
		log.Fatalf("Failed to execute template for %s: %v", outputPath, err)
	}

	fmt.Printf("Generated %s\n", outputFile)
}

// generateRedirect creates a simple HTML redirect page
func (g *GitHubPagesGenerator) generateRedirect(outputPath, target string) {
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
	outputFile := filepath.Join(g.OutputDir, outputPath)
	err := os.WriteFile(outputFile, []byte(redirectHTML), 0644)
	if err != nil {
		log.Fatalf("Failed to create redirect file %s: %v", outputFile, err)
	}

	fmt.Printf("Generated redirect from %s to %s\n", outputFile, target)
}

// copyDirectory recursively copies a directory tree
func (g *GitHubPagesGenerator) copyDirectory(src string, dst string) error {
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
			if err := g.copyDirectory(srcPath, dstPath); err != nil {
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

// Generate GitHub Pages version
func GenerateGitHubPages() {
	generator := NewGitHubPagesGenerator()
	generator.Run()
}
