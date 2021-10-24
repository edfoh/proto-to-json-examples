## Proto to Json Example

This is an example on how you can convert protobuf to JSON dynamically with [Go Templates](https://pkg.go.dev/text/template). This example looks into how to convert a proto message with a [oneOf](https://developers.google.com/protocol-buffers/docs/proto3#oneof). 

This is how it works.
1. convert protobuf to a map. 
2. Add go templates (json) that references the fields.
3. Run the matching go template with the converted map.

## Running the example

`make proto2json-templates` to run the example. Also run `make help` to see other useful commands.

## How it works

### Conversion

The [code](./internal/protoconvert/extract.go) takes in a `proto.Message` and converts that to a map. The [example.proto](./examples/example.proto) has a field `privileges` which is a `oneOf` for `Discount` or `FreeGift` message type. You would also pass in a `discriminator field name` so the conversion can retrieve the actual `oneOf` type specified in the payload, in this case the discriminator is `privileges`. The result of the conversion is the actual converted map and the extracted field name in the `oneOf`. 

### Go Templates

The templates are located in this [folder](./cmd/proto2json-templates/templates). The templates are loaded into the [embed fs](https://pkg.go.dev/embed) in [templates.go](./cmd/proto2json-templates/templates/templates.go). [base.tmpl](./cmd/proto2json-templates/customer/base.tmpl) is the base json structure, [discount.tmpl](./cmd/proto2json-templates/customer/discount.tmpl) and [free_gift.tmpl](./cmd/proto2json-templates/customer/free_gift.tmpl) are the json structures representing the `oneOf` types.

### Executing the templates

See [main.go](./cmd/proto2json-templates/main.go) to see how everything is hooked up.


