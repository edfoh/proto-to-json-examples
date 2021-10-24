package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"

	"github.com/edfoh/proto-to-json-poc/cmd/proto2json-templates/templates"
	"github.com/edfoh/proto-to-json-poc/internal/protobuf"
	"github.com/edfoh/proto-to-json-poc/internal/protoconvert"
	"google.golang.org/protobuf/proto"
)

func main() {
	convertCustomer(protobuf.CustomerWithDiscount(), "discount")
	fmt.Println()
	fmt.Println()
	convertCustomer(protobuf.CustomerWithFreeGift(), "free_gift")
}

func convertCustomer(customer proto.Message, oneOfName string) {
	fmt.Printf("converting customer data with %s", oneOfName)

	res, err := protoconvert.ToMap(customer, "privileges")
	if err != nil {
		fmt.Printf("error executing template: %v", err)
		panic(err)
	}

	t, err := parseTemplate("customer", oneOfName)
	if err != nil {
		fmt.Printf("error parsing template: %v", err)
		panic(err)
	}

	err = t.ExecuteTemplate(os.Stdout, "base", res.Out)
	if err != nil {
		fmt.Printf("error executing template: %v", err)
		panic(err)
	}
}

func parseTemplate(tmplType string, oneOfName string) (*template.Template, error) {
	name := fmt.Sprintf("%s-%s", tmplType, oneOfName)
	filePatterns := []string{
		fmt.Sprintf("%s/base.tmpl", tmplType),
		fmt.Sprintf("%s/%s.tmpl", tmplType, oneOfName),
	}
	return template.New(name).Funcs(template.FuncMap{
		"marshalMapToJSON": marshalMapToJSON,
	}).ParseFS(templates.CustomerTemplateFiles, filePatterns...)
}

func marshalMapToJSON(m map[string]interface{}) string {
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("error marshalMapToJSON: %v", err)
		return ""
	}
	return string(b)
}
