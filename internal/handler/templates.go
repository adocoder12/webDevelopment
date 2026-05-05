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
	Mutex       sync.RWMutex
	Dev         bool
	TemplateDir string
}

// TemplateData is the standard envelope passed to every template.
// Handler-specific data (Flash, Error, page content) goes in Content.
type TemplateData struct {
	IsAuthenticated bool
	UserName        string
	Flash           string
	Error           string
	Content         any
}

// NewTemplateRenderer initializes the specialist.
// It pre-parses all templates if not in Dev mode.
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

// newTemplateData builds the base data envelope for every render,
// automatically injecting auth state from the session.
func (app *Application) newTemplateData(r *http.Request) *TemplateData {
	td := &TemplateData{
		IsAuthenticated: app.isAuthenticated(r),
		Flash:           app.getFlash(r),
	}

	// If logged in, fetch the username from the session
	if td.IsAuthenticated {
		td.UserName = app.SessionManager.GetString(r.Context(), "userName")
	}

	return td
}

// renderTemplate executes a named template with the standard data envelope.
// Pass handler-specific data via td.Content, td.Error, etc. after calling
// newTemplateData.
func (app *Application) renderTemplate(w http.ResponseWriter, r *http.Request, filename string, td *TemplateData) {
	// Always start from a fresh base envelope so auth state is never stale
	if td == nil {
		td = app.newTemplateData(r)
	} else {
		// Inject auth state into the caller-supplied envelope
		td.IsAuthenticated = app.isAuthenticated(r)
		td.UserName = app.SessionManager.GetString(r.Context(), "userName")
		if td.Flash == "" {
			td.Flash = app.getFlash(r)
		}
	}

	var tmpl *template.Template
	var ok bool
	var err error

	app.Renderer.Mutex.RLock()
	tmpl, ok = app.Renderer.Cache[filename]
	app.Renderer.Mutex.RUnlock()

	if !ok || app.Renderer.Dev {
		tmpl, err = app.Renderer.ParseTemplate(filename)
		if err != nil {
			app.serverError(w, fmt.Errorf("could not parse template %s: %w", filename, err))
			return
		}

		if !app.Renderer.Dev {
			app.Renderer.Mutex.Lock()
			if _, exists := app.Renderer.Cache[filename]; !exists {
				app.Renderer.Cache[filename] = tmpl
			}
			app.Renderer.Mutex.Unlock()
		}
	}

	if err = tmpl.ExecuteTemplate(w, "base", td); err != nil {
		app.serverError(w, fmt.Errorf("could not execute template %s: %w", filename, err))
	}
}

func (tr *TemplateRenderer) ParseTemplate(filename string) (*template.Template, error) {
	pagePath := filepath.Join(tr.TemplateDir, filename)

	ts, err := template.ParseFiles(pagePath)
	if err != nil {
		return nil, fmt.Errorf("could not parse page %s: %w", filename, err)
	}

	layoutPattern := filepath.Join(tr.TemplateDir, "layouts/*.layout.html")
	matches, _ := filepath.Glob(layoutPattern)
	if len(matches) > 0 {
		ts, err = ts.ParseGlob(layoutPattern)
		if err != nil {
			return nil, fmt.Errorf("could not parse layouts: %w", err)
		}
	}

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
