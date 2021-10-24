package templates

import (
	"embed"
	"encoding/json"
	"fmt"
	"text/template"
)

//go:embed customer
var customerTemplateFiles embed.FS

func Parse(tmplType string, oneOfName string) (*template.Template, error) {
	name := fmt.Sprintf("%s-%s", tmplType, oneOfName)
	filePatterns := []string{
		fmt.Sprintf("%s/base.tmpl", tmplType),
		fmt.Sprintf("%s/%s.tmpl", tmplType, oneOfName),
	}
	return template.New(name).Funcs(template.FuncMap{
		"marshalMapToJSON": marshalMapToJSON,
	}).ParseFS(customerTemplateFiles, filePatterns...)
}

func marshalMapToJSON(m map[string]interface{}) string {
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("error marshalMapToJSON: %v", err)
		return ""
	}
	return string(b)
}
