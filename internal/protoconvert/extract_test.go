package protoconvert_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/edfoh/proto-to-json-poc/internal/protobuf"
	"github.com/edfoh/proto-to-json-poc/internal/protoconvert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToMap(t *testing.T) {
	expected := map[string]interface{}{
		"id":   int64(1),
		"name": "david",
		"attributes": map[string]interface{}{
			"phone":   "(415) 555-1212",
			"card_id": float64(1234),
			"tags":    []interface{}{"foo", "bar"},
		},
		"address": map[string]interface{}{
			"lines": []interface{}{
				map[string]interface{}{"value": "line1"},
				map[string]interface{}{"value": "line2"},
			},
		},
		"role": "normal",
		"extras": []interface{}{
			"1", "2",
		},
		"discount": map[string]interface{}{
			"max_discount": int64(50),
		},
	}

	in := protobuf.CustomerWithDiscount()
	res, err := protoconvert.ToMap(in, "privileges")
	require.NoError(t, err)

	js, _ := json.Marshal(res.Out)
	fmt.Printf(string(js))

	assert.Equal(t, "discount", res.DiscriminatorFieldName)
	assert.Equal(t, expected, res.Out)
}
