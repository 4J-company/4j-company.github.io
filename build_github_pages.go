package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"regexp"
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

			// Fix base URL tag which might not be rendering correctly
			fileContent = strings.Replace(fileContent,
				`<meta charset="UTF-8">
    <base href="{{ .BaseURL }}">`,
				`<meta charset="UTF-8">
    <base href="/">`,
				-1)

			// Replace the switchLanguage function with a completely rewritten version for GitHub Pages
			switchFunctionPattern := `function switchLanguage\(\) \{[\s\S]*?window\.location\.href = [\s\S]*?\}`
			newSwitchFunction := `function switchLanguage() {
    // Get the current path and check if we're on a Russian page
    const path = window.location.pathname;
    const isRussian = path.includes('_ru.html');
    
    // Determine the new path
    let newPath;
    if (isRussian) {
        // Switch from Russian to English
        newPath = path.replace('_ru.html', '.html');
    } else {
        // Switch from English to Russian
        newPath = path.replace('.html', '_ru.html');
    }
    
    // Handle special case for index page
    if (path === '/' || path === '/index.html') {
        newPath = '/index_ru.html';
    } else if (path === '/index_ru.html') {
        newPath = '/index.html';
    }
    
    // Handle case for paths ending with slash
    if (path.endsWith('/')) {
        if (isRussian) {
            newPath = path + 'index.html';
        } else {
            newPath = path + 'index_ru.html';
        }
    }
    
    console.log('Switching language from', path, 'to', newPath);
    window.location.href = newPath;
}`

			// Use regex to find and replace the switchLanguage function
			re := regexp.MustCompile(switchFunctionPattern)
			fileContent = re.ReplaceAllString(fileContent, newSwitchFunction)

			// Clean up any duplicated debug messages
			fileContent = strings.Replace(fileContent,
				`">
    // Debug message for language switching
    console.log("Language switching is set up. Current path:", window.location.pathname);
</script>`,
				`"></script>`,
				-1)

			// Replace dynamic links with static ones for both language versions
			if isRussian {
				// For Russian pages, ensure links point to Russian versions
				fileContent = replaceLink(fileContent, `href="/"`, `href="/index.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/?lang=ru{{else}}/{{end}}"`, `href="/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/features?lang=ru{{else}}/features{{end}}"`, `href="/features/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/examples?lang=ru{{else}}/examples{{end}}"`, `href="/examples/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/docs?lang=ru{{else}}/docs{{end}}"`, `href="/docs/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/download?lang=ru{{else}}/download{{end}}"`, `href="/download/index_ru.html"`)

				// Fix subproject links
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-graphics?lang=ru{{else}}/subprojects/mr-graphics{{end}}"`, `href="/subprojects/mr-graphics/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-math?lang=ru{{else}}/subprojects/mr-math{{end}}"`, `href="/subprojects/mr-math/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-importer?lang=ru{{else}}/subprojects/mr-importer{{end}}"`, `href="/subprojects/mr-importer/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-contractor?lang=ru{{else}}/subprojects/mr-contractor{{end}}"`, `href="/subprojects/mr-contractor/index_ru.html"`)
			} else {
				// For English pages, ensure links point to English versions
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/?lang=ru{{else}}/{{end}}"`, `href="/index.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/features?lang=ru{{else}}/features{{end}}"`, `href="/features/index.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/examples?lang=ru{{else}}/examples{{end}}"`, `href="/examples/index.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/docs?lang=ru{{else}}/docs{{end}}"`, `href="/docs/index.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/download?lang=ru{{else}}/download{{end}}"`, `href="/download/index.html"`)

				// Fix subproject links
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-graphics?lang=ru{{else}}/subprojects/mr-graphics{{end}}"`, `href="/subprojects/mr-graphics/index.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-math?lang=ru{{else}}/subprojects/mr-math{{end}}"`, `href="/subprojects/mr-math/index.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-importer?lang=ru{{else}}/subprojects/mr-importer{{end}}"`, `href="/subprojects/mr-importer/index.html"`)
				fileContent = replaceLink(fileContent, `href="{{if eq .Lang "ru"}}/subprojects/mr-contractor?lang=ru{{else}}/subprojects/mr-contractor{{end}}"`, `href="/subprojects/mr-contractor/index.html"`)
			}

			// Fix any direct links that might have been added during the page generation process
			// These would be links with no templates that need to be fixed for GitHub Pages
			if isRussian {
				// Fix standard links that have ?lang=ru format to use index_ru.html format
				fileContent = replaceLink(fileContent, `href="/?lang=ru"`, `href="/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="/features?lang=ru"`, `href="/features/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="/examples?lang=ru"`, `href="/examples/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="/docs?lang=ru"`, `href="/docs/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="/download?lang=ru"`, `href="/download/index_ru.html"`)

				// Fix subproject links
				fileContent = replaceLink(fileContent, `href="/subprojects/mr-graphics?lang=ru"`, `href="/subprojects/mr-graphics/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="/subprojects/mr-math?lang=ru"`, `href="/subprojects/mr-math/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="/subprojects/mr-importer?lang=ru"`, `href="/subprojects/mr-importer/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="/subprojects/mr-contractor?lang=ru"`, `href="/subprojects/mr-contractor/index_ru.html"`)

				fileContent = replaceLink(fileContent, `href="/examples"`, `href="/examples/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="/features"`, `href="/features/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="/docs"`, `href="/docs/index_ru.html"`)
				fileContent = replaceLink(fileContent, `href="/download"`, `href="/download/index_ru.html"`)
			} else {
				fileContent = replaceLink(fileContent, `href="/examples"`, `href="/examples/index.html"`)
				fileContent = replaceLink(fileContent, `href="/features"`, `href="/features/index.html"`)
				fileContent = replaceLink(fileContent, `href="/docs"`, `href="/docs/index.html"`)
				fileContent = replaceLink(fileContent, `href="/download"`, `href="/download/index.html"`)
			}

			// Fix language switch button visibility
			// Make sure it displays the correct flag and text
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

			// Add a prominent debug message at the top of the page
			fileContent = strings.Replace(fileContent,
				`<body class="bg-white text-black">`,
				`<body class="bg-white text-black">
    <script>
        console.log("Page language: `+getLanguageDisplay(isRussian)+`");
        console.log("Current path:", window.location.pathname);
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

// Helper function for language display
func getLanguageDisplay(isRussian bool) string {
	if isRussian {
		return "Russian"
	}
	return "English"
}

// Helper function to replace links without duplicating code
func replaceLink(content, oldLink, newLink string) string {
	return strings.Replace(content, oldLink, newLink, -1)
}

// Generate GitHub Pages version
func GenerateGitHubPages() {
	generator := NewGitHubPagesGenerator()
	generator.Run()
}
