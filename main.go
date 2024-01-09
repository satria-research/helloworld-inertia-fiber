package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/satria-research/inertia-fiber"
)

func main() {
	var optDir string
	flag.StringVar(&optDir, "dir", "", "project directory")
	flag.Parse()

	if optDir == "" {
		optDir, _ = os.Getwd()
	}

	e := fiber.New()

	e.Use(recover.New())
	e.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))

	// setup inertia
	r := inertia.NewRenderer()
	r.MustParseGlob(filepath.Join(optDir, "views/*.html"))
	r.ViteBasePath = "/dist/"
	r.AddViteEntryPoint("js/app.tsx")
	r.MustParseViteManifestFile(filepath.Join(optDir, "public/dist/manifest.json"))

	e.Use(inertia.Middleware(r))
	// e.Use(inertia.CSRF())

	e.Static("/", filepath.Join(optDir, "public"))

	e.Get("/", func(c *fiber.Ctx) error {
		c.Locals("Inertia", "Inertia")
		return inertia.Render(c, http.StatusOK, "Index", map[string]interface{}{
			"title":   "Hello, World! powered by inertia-fiber",
			"message": "Hello, World!",
		})
	})

	log.Fatal((e.Listen(":8080")))
}
