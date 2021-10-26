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
	fmt.Print("converting customer data with discount.\n")
	convertCustomer(protobuf.CustomerWithDiscount())

	fmt.Print("\n\n")

	fmt.Print("converting customer data with free gift - coupon.\n")
	convertCustomer(protobuf.CustomerWithFreeGiftCoupon())

	fmt.Print("\n\n")

	fmt.Print("converting customer data with free gift - item.\n")
	convertCustomer(protobuf.CustomerWithFreeGiftItem())
}

func convertCustomer(customer proto.Message) {

	res, err := protoconvert.ToMap(customer, "privileges", "free_gift.gift")
	if err != nil {
		fmt.Printf("error executing template: %v", err)
		panic(err)
	}

	t, err := templates.Parse("customer", res.OneOfFieldNames...)
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
