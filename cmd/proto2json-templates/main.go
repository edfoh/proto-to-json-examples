package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/edfoh/proto-to-json-examples/internal/customer"
	"github.com/edfoh/proto-to-json-examples/internal/protobuf"
)

func main() {
	fmt.Print("converting customer data with discount.\n")
	printJSON(customer.Convert(protobuf.CustomerWithDiscount()))

	fmt.Print("converting customer data with free gift - coupon.\n")
	printJSON(customer.Convert(protobuf.CustomerWithFreeGiftCoupon()))

	fmt.Print("converting customer data with free gift - item.\n")
	printJSON(customer.Convert(protobuf.CustomerWithFreeGiftItem()))
}

func printJSON(s string) {
	var b bytes.Buffer
	if err := json.Indent(&b, []byte(s), "", "  "); err != nil {
		panic(fmt.Sprintf("error with printing customer: %v", err))
	}
	fmt.Println(b.String())
	fmt.Println()
}
