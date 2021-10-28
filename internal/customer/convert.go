package customer

import (
	"bytes"
	"fmt"

	"github.com/edfoh/proto-to-json-examples/cmd/proto2json-templates/templates"
	"github.com/edfoh/proto-to-json-examples/internal/protoconvert"
	"google.golang.org/protobuf/proto"
)

func Convert(cust proto.Message) string {
	res, err := protoconvert.ToMap(cust, "privileges", "free_gift.gift")
	if err != nil {
		fmt.Printf("error executing template: %v", err)
		panic(err)
	}

	t, err := templates.Parse("customer", res.OneOfFieldNames...)
	if err != nil {
		fmt.Printf("error parsing template: %v", err)
		panic(err)
	}

	buf := new(bytes.Buffer)
	err = t.ExecuteTemplate(buf, "base", res.Out)
	if err != nil {
		fmt.Printf("error executing template: %v", err)
		panic(err)
	}
	return buf.String()
}
