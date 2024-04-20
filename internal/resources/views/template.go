package views

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/hako/durafmt"
)

const (
	base = "base.html"

	STATUS  = "status"
	SERVICE = "service"
)

var (
	templates      map[string]*template.Template
	ErrUnknownView = errors.New("view not found")
)

// Build all the templates at the once
func MustCompile() {
	templates = map[string]*template.Template{
		STATUS: template.Must(template.ParseFS(
			Views, base, "status.html",
		)),
		SERVICE: template.Must(template.ParseFS(
			Views, base, "service.html",
		)),
	}
}

// Execute a template and get the string data back
func Build(name string, data any) (string, error) {
	tmpl, ok := templates[name]
	if !ok {
		return "", ErrUnknownView
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type Builder struct{}

func (b Builder) Sprintf(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func (b Builder) PrettyDuration(duration time.Duration) string {
	return durafmt.Parse(duration).String()
}
