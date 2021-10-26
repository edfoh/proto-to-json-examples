package protobuf

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func CustomerWithDiscount() proto.Message {
	c := customer()
	c.Privileges = &Customer_Discount{
		Discount: &Discount{
			MaxDiscount: wrapperspb.Int64(50),
		},
	}
	return c
}

func CustomerWithFreeGiftCoupon() proto.Message {
	c := customer()
	c.Privileges = &Customer_FreeGift{
		FreeGift: &FreeGift{
			Gift: &FreeGift_Coupon{
				Coupon: &Coupon{
					Value: wrapperspb.Float(100),
				},
			},
		},
	}
	return c
}

func CustomerWithFreeGiftItem() proto.Message {
	c := customer()
	c.Privileges = &Customer_FreeGift{
		FreeGift: &FreeGift{
			Gift: &FreeGift_Item{
				Item: &Item{
					ItemCode: wrapperspb.String("ABC123"),
				},
			},
		},
	}
	return c
}

func customer() *Customer {
	props, _ := structpb.NewStruct(map[string]interface{}{
		"phone":   "(415) 555-1212",
		"card_id": 1234,
		"tags":    []interface{}{"foo", "bar"},
	})

	return &Customer{
		Id:         wrapperspb.Int64(int64(1)),
		Name:       wrapperspb.String(string("david")),
		Attributes: props,
		Address: &Address{
			Lines: []*AddressLine{
				{
					Value: wrapperspb.String("line1"),
				}, {
					Value: wrapperspb.String("line2"),
				},
			},
		},
		Role: Role_NORMAL,
		Extras: []*wrapperspb.StringValue{
			wrapperspb.String(string("1")),
			wrapperspb.String(string("2")),
		},
	}
}
