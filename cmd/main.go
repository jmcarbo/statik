package main

import (
	"flag"
	"fmt"
	"github.com/jmcarbo/statik/internal/token"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

var (
	ADir = flag.String("dir", "public", "directory to serve")
)

func main() {
	// Configuration variable for the static directory
	// Set via the STATIC_DIR environment variable or default to "./public"
	staticDir := "./public"
	if dir := os.Getenv("STATIC_DIR"); dir != "" {
		staticDir = dir
	}
	flag.Parse()
	if *ADir != "" {
		staticDir = *ADir
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

				auth := c.Get("Authorization")
				// Data to pass to the template, including headers
				data := fiber.Map{
					"Headers": c.GetReqHeaders(),
					"Method":  c.Method(),
					"Path":    c.Path(),
					"Query":   c.OriginalURL(),
					"Token":   token.GetToken(auth),
				}

				// Render the template
				b, err := os.ReadFile(filePath)
				if err != nil {
					log.Printf("Error reading file %s: %v", filePath, err)
					return fiber.ErrInternalServerError
				}
				tmpl, err := template.New("test").Parse(string(b))
				if err != nil {
					log.Printf("Error parsing template %s: %v", filePath, err)
					return fiber.ErrInternalServerError
				}
				buf := strings.Builder{}
				err = tmpl.Execute(&buf, data)
				if err != nil {
					log.Printf("Error executing template %s: %v", filePath, err)
					return fiber.ErrInternalServerError
				}

				c.Set("Content-Type", "text/html")
				return c.SendString(buf.String())
			} else {
				// Serve static file
				return c.SendFile("./"+filePath, true)
			}
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