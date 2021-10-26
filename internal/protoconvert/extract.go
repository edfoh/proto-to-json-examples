package protoconvert

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ProtoToMapResult struct {
	OneOfFieldNames []string
	Out             map[string]interface{}
}

// ToMap takes in a proto.Message and converts that to a map.
//
// if oneOfPaths is specified, it will inspect the proto message for the actual discriminator field name for
// in the oneOf specification. The path should be 'event' or 'some_event.another' for nested oneOf
//
// returns a ProtoToMapResult
func ToMap(in proto.Message, oneOfPaths ...string) (*ProtoToMapResult, error) {
	oneOfFieldNames, err := extractOneOfPaths(in, oneOfPaths...)
	if err != nil {
		return nil, err
	}

	return &ProtoToMapResult{
		OneOfFieldNames: oneOfFieldNames,
		Out:             extractProtoMessage(in),
	}, nil
}

func extractOneOfPaths(in proto.Message, oneOfPaths ...string) ([]string, error) {
	var res []string
	if oneOfPaths == nil || len(oneOfPaths) == 0 {
		return res, nil
	}

	message := in.ProtoReflect()
	for _, path := range oneOfPaths {
		parts := strings.Split(path, ".")

		if len(parts) == 1 {
			out, err := extractOneOfFieldName(message, parts[0])
			if err != nil {
				return res, nil
			}
			res = append(res, out)
			continue
		}

		fd := message.Descriptor().Fields().ByName(protoreflect.Name(parts[0]))
		if fd.Kind() == protoreflect.MessageKind {
			nextPath := strings.Join(parts[1:], ".")
			value := message.Get(fd)
			out, err := extractOneOfPaths(value.Message().Interface(), nextPath)
			if err != nil {
				return res, err
			}
			res = append(res, out...)
		}
	}

	return res, nil
}

func extractOneOfFieldName(message protoreflect.Message, name string) (string, error) {
	oneOfDesc := message.Descriptor().Oneofs().ByName(protoreflect.Name(name))
	if oneOfDesc == nil {
		return "", fmt.Errorf("unable to find oneOf field name '%s'", name)
	}

	oneOfFieldName := message.WhichOneof(oneOfDesc)
	if oneOfFieldName == nil {
		return "", fmt.Errorf("oneOf field name '%s' has no value populated", name)
	}

	return string(oneOfFieldName.Name()), nil
}

func extractProtoMessage(in proto.Message) map[string]interface{} {
	res := map[string]interface{}{}

	msg := in.ProtoReflect()
	msg.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		fieldValue := msg.Get(fd)
		if v := extractValue(fd, fieldValue); isNotEmpty(v) {
			res[string(fd.Name())] = v
		}
		return true
	})

	return res
}

func extractValue(fd protoreflect.FieldDescriptor, v protoreflect.Value) interface{} {
	if fd.IsList() {
		return extractList(v.List())
	}
	if fd.IsMap() {
		return extractMap(v.Map(), fd.MapValue())
	}
	if fd.Kind() == protoreflect.MessageKind {
		return extractValueFromMessage(v.Message().Interface())
	}
	if fd.Kind() == protoreflect.EnumKind {
		enumString := string(fd.Enum().Values().ByNumber(v.Enum()).Name())
		return strings.ToLower(enumString)
	}
	return v.Interface()
}

func extractList(l protoreflect.List) interface{} {
	var res []interface{}
	for i := 0; i < l.Len(); i++ {
		res = append(res, extractValueFromMessage(l.Get(i).Message().Interface()))
	}
	return res
}

func extractMap(m protoreflect.Map, valueFd protoreflect.FieldDescriptor) interface{} {
	res := map[string]interface{}{}

	m.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
		res[k.String()] = extractValue(valueFd, v)
		return true
	})

	return res
}

func extractValueFromMessage(m proto.Message) interface{} {
	if m == nil {
		return nil
	}

	switch in := m.(type) {
	case nil:
		return nil
	case *wrapperspb.DoubleValue:
		return in.GetValue()
	case *wrapperspb.FloatValue:
		return in.GetValue()
	case *wrapperspb.Int32Value:
		return in.GetValue()
	case *wrapperspb.Int64Value:
		return in.GetValue()
	case *wrapperspb.StringValue:
		return in.GetValue()
	case *wrapperspb.BoolValue:
		return in.GetValue()
	case *structpb.Struct:
		return extractStruct(in)
	default:
		return extractProtoMessage(m)
	}
}

func extractStruct(s *structpb.Struct) interface{} {
	if s == nil {
		return nil
	}

	res := map[string]interface{}{}
	for k, v := range s.Fields {
		res[k] = v.AsInterface()
	}
	return res
}

func isNotEmpty(v interface{}) bool {
	if m, ok := v.(map[string]interface{}); ok {
		return len(m) > 0
	}
	return v != nil
}
