// pkg/utils/templates.go
package utils

import (
	"fmt"
	"html/template"
	"strings"
)

type Message struct {
	Content string
	Type    string
}

func FirstFlashOrEmpty(flashes []interface{}) string {
	if len(flashes) > 0 {
		return flashes[0].(string)
	}
	return ""
}

func Add(a, b int) int { return a + b }

func Sub(a, b int) int { return a - b }

// Función para generar HTML de las hojas de estilo (versión mejorada)
func GenerateStylesHTML(baseURL string, styles []string) template.HTML {
	if baseURL == "" {
		baseURL = "/static/"
	} else if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	var html strings.Builder
	for _, style := range styles {
		style = strings.Trim(style, "/")
		if !strings.HasSuffix(style, ".css") {
			style += ".css"
		}
		html.WriteString(fmt.Sprintf(`<link rel="stylesheet" href="%s%s">`, baseURL, style))
	}

	return template.HTML(html.String())
}

// Función para generar HTML de los scripts (versión mejorada)
func GenerateScriptsHTML(baseURL string, scripts []string) template.HTML {
	if baseURL == "" {
		baseURL = "/static/"
	} else if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	var html strings.Builder
	for _, script := range scripts {
		script = strings.Trim(script, "/")
		if !strings.HasSuffix(script, ".js") {
			script += ".js"
		}
		html.WriteString(fmt.Sprintf(`<script src="%s%s"></script>`, baseURL, script))
	}

	return template.HTML(html.String())
}
