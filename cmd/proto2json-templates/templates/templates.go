package templates

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
)

//go:embed customer
var customerTemplateFiles embed.FS

func Parse(tmplType string, oneOfNames ...string) (*template.Template, error) {
	name := fmt.Sprintf("%s-%s", tmplType, strings.Join(oneOfNames, "-"))
	filePatterns := []string{
		fmt.Sprintf("%s/base.tmpl", tmplType),
	}
	for _, name := range oneOfNames {
		filePatterns = append(filePatterns, fmt.Sprintf("%s/%s.tmpl", tmplType, name))
	}

	return template.New(name).Funcs(template.FuncMap{
		"marshalMapToJSON": marshalMapToJSON,
	}).ParseFS(customerTemplateFiles, filePatterns...)
}

func marshalMapToJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("error marshalMapToJSON: %v", err)
		return ""
	}
	return string(b)
}
