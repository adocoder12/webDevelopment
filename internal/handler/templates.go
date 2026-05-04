package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

type TemplateRenderer struct {
	Cache       map[string]*template.Template
	Mutex       sync.RWMutex // Use RWMutex for better performance
	Dev         bool
	TemplateDir string
}

// NewTemplateRenderer initializes the specialist.
// It pre-parses if not in Dev mode.
func NewTemplateRenderer(dir string, isDev bool) (*TemplateRenderer, error) {
	renderer := &TemplateRenderer{
		Cache:       make(map[string]*template.Template),
		Dev:         isDev,
		TemplateDir: dir,
	}

	if !isDev {
		pages, err := filepath.Glob(filepath.Join(dir, "*.html"))
		if err != nil {
			return nil, err
		}

		for _, page := range pages {
			name := filepath.Base(page)
			ts, err := template.ParseFiles(page)
			if err != nil {
				return nil, err
			}
			renderer.Cache[name] = ts
		}
	}

	return renderer, nil
}

// renderTemplate now uses the Renderer specialist
func (app *Application) renderTemplate(w http.ResponseWriter, filename string, data any) {
	var tmpl *template.Template
	var ok bool
	var err error // Pre-declare err to avoid scope issues

	// 1. Thread-safe check of the cache
	app.Renderer.Mutex.RLock()
	tmpl, ok = app.Renderer.Cache[filename]
	app.Renderer.Mutex.RUnlock()

	// 2. Logic for missing cache or Dev Mode
	if !ok || app.Renderer.Dev {

		// Use '=' here, not ':=' to update the variable declared at the top
		tmpl, err = app.Renderer.ParseTemplate(filename)
		if err != nil {
			app.serverError(w, fmt.Errorf("could not parse template %s: %w", filename, err))
			return
		}

		// 3. Update cache if not in Dev mode
		if !app.Renderer.Dev {
			app.Renderer.Mutex.Lock()
			if _, exists := app.Renderer.Cache[filename]; !exists {
				app.Renderer.Cache[filename] = tmpl
			}
			app.Renderer.Mutex.Unlock()
		}
	}

	// 4. Execute the template
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, fmt.Errorf("could not execute template %s: %w", filename, err))
	}
}

func (tr *TemplateRenderer) ParseTemplate(filename string) (*template.Template, error) {
	pagePath := filepath.Join(tr.TemplateDir, filename)

	ts, err := template.ParseFiles(pagePath)
	if err != nil {
		return nil, fmt.Errorf("could not parse page %s: %w", filename, err)
	}

	// Check for layouts
	layoutPattern := filepath.Join(tr.TemplateDir, "layouts/*.layout.html")
	matches, _ := filepath.Glob(layoutPattern)
	if len(matches) > 0 {
		ts, err = ts.ParseGlob(layoutPattern)
		if err != nil {
			return nil, fmt.Errorf("could not parse layouts: %w", err)
		}
	}

	// Check for partials
	partialPattern := filepath.Join(tr.TemplateDir, "partials/*.partial.html")
	matches, _ = filepath.Glob(partialPattern)
	if len(matches) > 0 {
		ts, err = ts.ParseGlob(partialPattern)
		if err != nil {
			return nil, fmt.Errorf("could not parse partials: %w", err)
		}
	}

	return ts, nil
}
