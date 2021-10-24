package main

import (
	"fmt"
	"os"

	"github.com/edfoh/proto-to-json-examples/cmd/proto2json-templates/templates"
	"github.com/edfoh/proto-to-json-examples/internal/protobuf"
	"github.com/edfoh/proto-to-json-examples/internal/protoconvert"
	"google.golang.org/protobuf/proto"
)

func main() {
	convertCustomer(protobuf.CustomerWithDiscount(), "discount")
	fmt.Print("\n\n")
	convertCustomer(protobuf.CustomerWithFreeGift(), "free_gift")
}

func convertCustomer(customer proto.Message, oneOfName string) {
	fmt.Printf("converting customer data with %s", oneOfName)

	res, err := protoconvert.ToMap(customer, "privileges")
	if err != nil {
		fmt.Printf("error executing template: %v", err)
		panic(err)
	}

	t, err := templates.Parse("customer", oneOfName)
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
