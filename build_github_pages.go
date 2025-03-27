package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	g.generateRedirect("docs/index_ru.html", "/index_ru.html")
	g.generateRedirect("download/index.html", "/")
	g.generateRedirect("download/index_ru.html", "/index_ru.html")

	// Fix all the HTML files to work with GitHub Pages static structure
	g.fixLanguageLinks()
}

// generatePage renders a template to a static HTML file
func (g *GitHubPagesGenerator) generatePage(outputPath, layoutPath, contentPath string, data PageData) {
	// Parse templates
	tmpl, err := template.ParseFiles(layoutPath, "templates/translations.html", contentPath)
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

	// Add Lang field and Year to data
	data.Lang = "en"
	data.Year = time.Now().Year()
	data.BaseURL = "/" // Base URL for GitHub Pages

	// Execute template and write to file
	if err := tmpl.ExecuteTemplate(file, "layout", data); err != nil {
		log.Fatalf("Failed to execute template for %s: %v", outputPath, err)
	}

	fmt.Printf("Generated %s\n", outputFile)

	// If we're generating the main English version, also generate Russian version
	if !strings.Contains(outputPath, "_ru") {
		// Create output file for Russian version
		ruOutputPath := strings.Replace(outputPath, "index.html", "index_ru.html", 1)
		ruOutputFile := filepath.Join(g.OutputDir, ruOutputPath)
		ruFile, err := os.Create(ruOutputFile)
		if err != nil {
			log.Fatalf("Failed to create file %s: %v", ruOutputFile, err)
		}
		defer ruFile.Close()

		// Set language to Russian for this version
		data.Lang = "ru"

		// Execute template and write to file
		if err := tmpl.ExecuteTemplate(ruFile, "layout", data); err != nil {
			log.Fatalf("Failed to execute template for %s: %v", ruOutputPath, err)
		}

		fmt.Printf("Generated %s\n", ruOutputFile)
	}
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

// fixLanguageLinks modifies the generated HTML files to make language switching work with static files
func (g *GitHubPagesGenerator) fixLanguageLinks() {
	fmt.Println("Fixing language links for GitHub Pages...")

	// Process all HTML files in the output directory
	err := filepath.Walk(g.OutputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".html" {
			// Read file
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %v", path, err)
			}

			fileContent := string(content)
			isRussian := strings.Contains(path, "_ru.html")

			// Replace the switchLanguage function with a more robust one for GitHub Pages
			fileContent = strings.Replace(fileContent,
				`function switchLanguage() {
            const currentUrl = new URL(window.location.href);
            const currentLang = currentUrl.searchParams.get('lang') || 'en';
            const newLang = currentLang === 'en' ? 'ru' : 'en';
            
            // Remove any existing lang parameter and add the new one
            currentUrl.searchParams.delete('lang');
            currentUrl.searchParams.set('lang', newLang);
            
            // Switch language without animation
            window.location.href = currentUrl.toString();
        }`,
				`function switchLanguage() {
            // For GitHub Pages, we use separate HTML files for different languages
            const path = window.location.pathname;
            const isRoot = path === "/" || path.endsWith("/");
            const isRussian = path.includes("_ru.html") || path.endsWith("_ru/");
            
            let newPath;
            if (isRussian) {
                // Switch from Russian to English
                if (isRoot) {
                    newPath = "/";
                } else {
                    newPath = path.replace("_ru.html", ".html").replace("_ru/", "/");
                }
            } else {
                // Switch from English to Russian
                if (isRoot) {
                    newPath = "/index_ru.html";
                } else if (path.endsWith("/")) {
                    newPath = path + "index_ru.html";
                } else if (path.endsWith(".html")) {
                    newPath = path.replace(".html", "_ru.html");
                } else {
                    // Path without .html extension
                    const lastSlashIndex = path.lastIndexOf("/");
                    if (lastSlashIndex !== -1) {
                        const basePath = path.substring(0, lastSlashIndex + 1);
                        const pageName = path.substring(lastSlashIndex + 1);
                        newPath = basePath + pageName + "_ru.html";
                    } else {
                        newPath = path + "_ru.html";
                    }
                }
            }
            window.location.href = newPath;
        }`,
				-1)

			// Replace dynamic links with static ones
			if isRussian {
				// For Russian pages, ensure links point to Russian versions
				fileContent = strings.Replace(fileContent, `href="/"`, `href="/index.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/?lang=ru{{else}}/{{end}}"`, `href="/index_ru.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/features?lang=ru{{else}}/features{{end}}"`, `href="/features/index_ru.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/examples?lang=ru{{else}}/examples{{end}}"`, `href="/examples/index_ru.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/docs?lang=ru{{else}}/docs{{end}}"`, `href="/docs/index_ru.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/download?lang=ru{{else}}/download{{end}}"`, `href="/download/index_ru.html"`, -1)

				// Fix subproject links
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-graphics?lang=ru{{else}}/subprojects/mr-graphics{{end}}"`, `href="/subprojects/mr-graphics/index_ru.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-math?lang=ru{{else}}/subprojects/mr-math{{end}}"`, `href="/subprojects/mr-math/index_ru.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-importer?lang=ru{{else}}/subprojects/mr-importer{{end}}"`, `href="/subprojects/mr-importer/index_ru.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-contractor?lang=ru{{else}}/subprojects/mr-contractor{{end}}"`, `href="/subprojects/mr-contractor/index_ru.html"`, -1)
			} else {
				// For English pages, ensure links point to English versions
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/?lang=ru{{else}}/{{end}}"`, `href="/index.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/features?lang=ru{{else}}/features{{end}}"`, `href="/features/index.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/examples?lang=ru{{else}}/examples{{end}}"`, `href="/examples/index.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/docs?lang=ru{{else}}/docs{{end}}"`, `href="/docs/index.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/download?lang=ru{{else}}/download{{end}}"`, `href="/download/index.html"`, -1)

				// Fix subproject links
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-graphics?lang=ru{{else}}/subprojects/mr-graphics{{end}}"`, `href="/subprojects/mr-graphics/index.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-math?lang=ru{{else}}/subprojects/mr-math{{end}}"`, `href="/subprojects/mr-math/index.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-importer?lang=ru{{else}}/subprojects/mr-importer{{end}}"`, `href="/subprojects/mr-importer/index.html"`, -1)
				fileContent = strings.Replace(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-contractor?lang=ru{{else}}/subprojects/mr-contractor{{end}}"`, `href="/subprojects/mr-contractor/index.html"`, -1)
			}

			// For both language versions, add base href to ensure absolute paths work
			// This is important for GitHub Pages, especially with custom domains or repository paths
			fileContent = strings.Replace(fileContent,
				`<meta charset="UTF-8">`,
				`<meta charset="UTF-8">
    <base href="{{ .BaseURL }}">`,
				-1)

			// Fix language switch button visibility
			// The issue could be with the conditional templates in layout.html not evaluating correctly
			if isRussian {
				fileContent = strings.Replace(fileContent,
					`{{if eq .Lang "en"}}🇷🇺{{else}}🇺🇸{{end}}`,
					`🇺🇸`,
					-1)
				fileContent = strings.Replace(fileContent,
					`{{template "lang.switch" .}}`,
					`English`,
					-1)
			} else {
				fileContent = strings.Replace(fileContent,
					`{{if eq .Lang "en"}}🇷🇺{{else}}🇺🇸{{end}}`,
					`🇷🇺`,
					-1)
				fileContent = strings.Replace(fileContent,
					`{{template "lang.switch" .}}`,
					`Русский`,
					-1)
			}

			// Add a console message for debugging purposes
			fileContent = strings.Replace(fileContent,
				`</script>`,
				`
    // Debug message for language switching
    console.log("Language switching is set up. Current path:", window.location.pathname);
</script>`,
				-1)

			// Write modified content back to file
			err = os.WriteFile(path, []byte(fileContent), 0644)
			if err != nil {
				return fmt.Errorf("error writing file %s: %v", path, err)
			}

			fmt.Printf("Fixed language links in %s\n", path)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error fixing language links: %v", err)
	}
}

// Generate GitHub Pages version
func GenerateGitHubPages() {
	generator := NewGitHubPagesGenerator()
	generator.Run()
}
