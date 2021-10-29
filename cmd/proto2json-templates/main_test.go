package main_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/edfoh/proto-to-json-examples/internal/customer"
	"github.com/edfoh/proto-to-json-examples/internal/protobuf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	protojsonDir = "testdata/protojson"
	jsonDir      = "testdata/json"
)

func TestConvert(t *testing.T) {
	testCustomer(t)
}

func testCustomer(t *testing.T) {
	files, err := ioutil.ReadDir(protojsonDir)
	require.NoError(t, err)

	for _, file := range files {
		filename := file.Name()
		filenameWithoutExt := fileNameWithoutExtension(filename)

		t.Run(fmt.Sprintf("testcase: protojson file %s", filename), func(t *testing.T) {
			c := loadCustomer(t, filename)

			actual := customer.Convert(c)
			expected := loadExpectedJsonOutput(t, filenameWithoutExt)

			assert.JSONEq(t, expected, actual)
		})
	}
}

func loadCustomer(t *testing.T, filename string) protoreflect.ProtoMessage {
	// load the protojson file
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", protojsonDir, filename))
	require.NoError(t, err)

	var c protobuf.Customer
	err = protojson.Unmarshal(b, &c)
	require.NoError(t, err)

	return &c
}

func loadExpectedJsonOutput(t *testing.T, filenameWithoutExt string) string {
	b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", jsonDir, filenameWithoutExt))
	require.NoError(t, err)

	return string(b)
}

func fileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
