# 4J model-renderer Website

A modern, responsive website for the 4J model-renderer engine and its modules, built with Go and featuring a clean black-and-white design aesthetic.

## Features

- Modern, high-contrast black-and-white design
- Responsive layout with mobile support
- Syntax highlighting for code examples
- HTMX for dynamic content loading
- Go backend with easy-to-maintain templates
- GitHub Pages deployment ready

## Live Demo

Visit the live site at: `https://4j-company.github.io/mr-website/`

## Modules Showcase

The website features detailed pages for each of the model-renderer modules:

- **mr-graphics**: High-performance Vulkan rendering library
- **mr-importer**: Asset importing and processing system
- **mr-contractor**: Declarative task execution library
- **mr-math**: Linear algebra and mathematics utilities

## Development

### Prerequisites

- Go 1.21 or later
- Git

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/4j-company/mr-website.git
cd mr-website
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run main.go
```

4. Open your browser and navigate to:
```
http://localhost:4747
```

### Project Structure

```
.
├── main.go                 # Main server file
├── go.mod                 # Go module file
├── templates/             # HTML templates
│   ├── layout.html        # Base layout template
│   ├── home.html          # Home page
│   ├── features.html      # Features overview
│   ├── examples.html      # Example showcase
│   └── subprojects/       # Module-specific pages
│       ├── mr-graphics.html
│       ├── mr-importer.html
│       ├── mr-contractor.html
│       └── mr-math.html
├── static_generator.go    # Static site generator for GitHub Pages
├── docs/                  # Generated static site (for GitHub Pages)
│   └── ...
└── .github/workflows/    # GitHub Actions workflow
    └── pages.yml         # Workflow for GitHub Pages deployment
```

## Deployment to GitHub Pages

### Manual Deployment

1. Generate the static site:
```bash
go run . --github-pages
```

2. Commit and push the changes:
```bash
git add docs
git commit -m "Update static site"
git push
```

3. Configure GitHub Pages in your repository settings to use the `/docs` folder on the master branch.

### Automated Deployment

The repository includes a GitHub Actions workflow that automatically builds and deploys the site to GitHub Pages whenever changes are pushed to the master branch.

## Customization

- Edit the HTML templates in the `templates/` directory to modify content
- Update styling in `templates/layout.html`
- Adjust the static site generator in `static_generator.go` if needed

## License

MIT License 
