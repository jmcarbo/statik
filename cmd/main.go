package main

import (
	"flag"
	"fmt"
	"github.com/jmcarbo/statik/internal/token"
	"html/template"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var (
	ADir        = flag.String("dir", "public", "directory to serve")
	SPA         = flag.Bool("spa", false, "single page application mode")
	MultiTenant = flag.Bool("multitenant", false, "multi tenant mode")
)

func init() {
	mimeTypes := map[string]string{
		".ts":   "application/typescript",
		".tsx":  "application/typescript",
		".d.ts": "application/typescript",
	}

	for ext, mimeType := range mimeTypes {
		if err := mime.AddExtensionType(ext, mimeType); err != nil {
			log.Fatalf("Error adding MIME type for %s: %v", ext, err)
		}
	}
}

func getTenantPath(path string) string {
	// Split the path into parts
	parts := strings.Split(path, "/")

	// If the path has at least 2 parts, the first part is the tenant
	if len(parts) >= 1 {
		return filepath.Join(parts[2:]...)
	}

	return path
}

func main() {
	// Configuration variable for the static directory
	// Set via the STATIC_DIR environment variable or default to "./public"
	staticDir := "./public"
	if dir := os.Getenv("STATIC_DIR"); dir != "" {
		staticDir = dir
	}
	flag.Parse()

	if *SPA {
		log.Println("Single Page Application mode enabled")
	}

	if *ADir != "" {
		staticDir = *ADir
	}

	if *MultiTenant {
		log.Println("Multi Tenant mode enabled")
	}

	log.Printf("Serving files from %s", staticDir)

	// Create a new HTML template engine with the static directory as the views folder
	engine := html.New(staticDir, ".html")

	// Initialize the Fiber app with the template engine
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Middleware to handle serving static files and rendering templates
	app.Use(func(c *fiber.Ctx) error {
		// Get the requested path
		path := c.Path()

		// Multi Tenant mode
		if *MultiTenant {
			path = getTenantPath(path)
		}

		// Construct the full file path
		filePath := filepath.Join(staticDir, path)
		log.Printf("Serving file %s", filePath)

		// Check if the file or directory exists
		if info, err := os.Stat(filePath); err == nil {
			if info.IsDir() {
				// If it's a directory, look for index.html
				filePath = filepath.Join(filePath, "index.html")
				if _, err := os.Stat(filePath); err != nil {
					return fiber.ErrNotFound
				}
				// Update the path to include index.html
				path = filepath.Join(path, "index.html")
			}

			// Check if the file is an HTML file
			if filepath.Ext(filePath) == ".html" {

				// Template path relative to the views directory
				templatePath, err := filepath.Rel(staticDir, filePath)

				if err != nil {
					return fiber.ErrInternalServerError
				}
				log.Printf("Template path %s", templatePath)

				auth := c.Get("Authentication")

				// Data to pass to the template, including headers
				data := fiber.Map{
					"Headers": c.GetReqHeaders(),
					"Method":  c.Method(),
					"Path":    c.Path(),
					"Query":   c.OriginalURL(),
					"Token":   token.GetToken(auth),
				}

				// Render the template
				tplRes, err := RenderTemplate(filePath, data)
				if err != nil {
					return err
				}

				c.Set("Content-Type", "text/html")
				return c.SendString(tplRes)
			} else {
				// Serve static file
				return c.SendFile("./"+filePath, true)
			}
		}

		log.Printf("File not found %s", filePath)
		if *SPA {
			// Serve index.html for SPA

			filePath = filepath.Join(staticDir, "index.html")
			if _, err := os.Stat(filePath); err != nil {
				return fiber.ErrNotFound
			}
			log.Printf("Serving index.html for SPA %s", filePath)
			auth := c.Get("Authentication")

			// Data to pass to the template, including headers
			data := fiber.Map{
				"Headers": c.GetReqHeaders(),
				"Method":  c.Method(),
				"Path":    c.Path(),
				"Query":   c.OriginalURL(),
				"Token":   token.GetToken(auth),
			}

			// Render the template
			tplRes, err := RenderTemplate(filePath, data)
			if err != nil {
				return err
			}

			c.Set("Content-Type", "text/html")
			return c.SendString(tplRes)
		}

		// File not found
		return fiber.ErrNotFound
	})

	// Start the server on port 3000
	err := app.Listen(":3000")
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func RenderTemplate(filePath string, data any) (string, error) {
	// Render the template
	b, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file %s: %v", filePath, err)
		return "", fiber.ErrInternalServerError
	}
	tmpl, err := template.New("test").Parse(string(b))
	if err != nil {
		log.Printf("Error parsing template %s: %v", filePath, err)
		return "", fiber.ErrInternalServerError
	}
	buf := strings.Builder{}
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Printf("Error executing template %s: %v", filePath, err)
		return "", fiber.ErrInternalServerError
	}
	return buf.String(), nil
}
