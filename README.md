## Proto to Json Example

This is an example on how you can convert protobuf to JSON dynamically with [Go Templates](https://pkg.go.dev/text/template). This example looks into how to convert a proto message with a [oneOf](https://developers.google.com/protocol-buffers/docs/proto3#oneof). 

This is how it works.
1. convert protobuf to a map. 
2. Add go templates (json) that references the fields.
3. Run the matching go template with the converted map.

## Running the example

`make proto2json-templates` to run the example. Also run `make help` to see other useful commands.

## How it works

We take this [example.proto](./examples/example.proto)

``` proto
enum Role {
	UNKNOWN = 0;
	NORMAL = 1;
	VIP = 2;
}

message Address {
    repeated AddressLine lines = 1;
}

message AddressLine {
    google.protobuf.StringValue value = 1;
}

message Discount {
    google.protobuf.Int64Value max_discount = 1;
}

message FreeGift {
    oneof gift {
		Coupon coupon = 1;
		Item item = 2;
	}
}

message Coupon {
	google.protobuf.FloatValue value = 1;
}

message Item {
	google.protobuf.StringValue item_code = 1;
}

message Customer {
	google.protobuf.Int64Value id = 1;
	google.protobuf.StringValue name = 2;
	google.protobuf.Struct attributes = 3;
	Address address = 4;
    Role role = 5;
	repeated google.protobuf.StringValue extras = 6;
	oneof privileges {
		Discount discount = 7;
		FreeGift free_gift = 8;
	}
}

```

and using go templates [base.tmpl](./cmd/proto2json-templates/templates/customer/base.tmpl) and [discount.tmpl](./cmd/proto2json-templates/templates/customer/discount.tmpl).

```
// base.tmpl
{{define "base"}}{
    "customer": {
        "id": "{{.id}}",
        "name": "{{.name}}",
        "address": {{marshalMapToJSON .address}},
        "misc": {
            "attributes": {{marshalMapToJSON .attributes}},
            "role": "{{.role}}",
            "extras": {{marshalMapToJSON .extras}}
        },
        "privileges": {
            {{template "oneOf1" .}}
        }
    }
}{{end}}

// discount.tmpl
{{define "oneOf1"}}
"discount": {
    "max_discount": {{.discount.max_discount}}
}{{end}}
```

we can get these JSON outputs
``` json
// customer with discount
{
  "customer": {
    "id": "1",
    "name": "david",
    "address": {
      "lines": [
        {
          "value": "line1"
        },
        {
          "value": "line2"
        }
      ]
    },
    "misc": {
      "attributes": {
        "card_id": 1234,
        "phone": "(415) 555-1212",
        "tags": [
          "foo",
          "bar"
        ]
      },
      "role": "normal",
      "extras": [
        "1",
        "2"
      ]
    },
    "privileges": {
      "discount": {
        "max_discount": 50
      }
    }
  }
}
```

This also supports nested `oneOf`.

### Conversion

The [code](./internal/protoconvert/extract.go) takes in a `proto.Message` and converts that to a map. The [example.proto](./examples/example.proto) has a field `privileges` which is a `oneOf` for `Discount` or `FreeGift` message type. You would also pass in a `oneOf path` so the conversion can retrieve the actual `oneOf` type specified in the payload, in this case the path is `privileges`. It also supports nested `oneOf path`, in this case `free_gift.gift`.

The result of the conversion is the actual converted map and the extracted oneOf field names. 

### Go Templates

The templates are located in this [folder](./cmd/proto2json-templates/templates). The templates are loaded into the [embed fs](https://pkg.go.dev/embed) in [templates.go](./cmd/proto2json-templates/templates/templates.go). [base.tmpl](./cmd/proto2json-templates/customer/base.tmpl) is the base json structure, [discount.tmpl](./cmd/proto2json-templates/customer/discount.tmpl) and [free_gift.tmpl](./cmd/proto2json-templates/customer/free_gift.tmpl) are the json structures representing the `oneOf` types. `free_gift.tmpl` also has a oneOf, which illustrates how you would template nested `oneOf`.
 
### Executing the templates

See [main.go](./cmd/proto2json-templates/main.go) to see how everything is hooked up.

### Testing

We are using [protojson](https://developers.google.com/protocol-buffers/docs/proto3#json) as source input files, which will go through the conversion. We will also have expected json output files that will be used to verify that they match the actual converted ones.

```
cmd/
├─ proto2json-templates/
│  ├─ testdata/
│  │  ├─ protojson/
│  │  │  ├─ customer.json
│  │  │  ├─ ...
│  │  ├─ json/
│  │  │  ├─ customer_discount.json
│  │  │  ├─ ...
|  |─ main_test.go

```

To add more tests to cover the different scenarios, adding input files in `testdata/protojson` folder and expected output in `testdata/json` folder. The filenames need to match to execute a test case. 

see [main_test.go](./cmd/proto2json-templates/main_test.go) for more details.

