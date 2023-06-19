package promec

import (
	"fmt"
	"github.com/civet148/log"
	"reflect"
)

func parseLabels(obj LabelObject) (labels, values []string) {

	if obj == nil {
		return
	}
	typ := reflect.TypeOf(obj)
	val := reflect.ValueOf(obj)

	for {
		if typ.Kind() != reflect.Ptr { // pointer type
			break
		}
		typ = typ.Elem()
		val = val.Elem()
	}

	kind := typ.Kind()
	switch kind {
	case reflect.Struct:
		{
			return parseStructFields(typ, val)
		}
	//case reflect.Slice:
	//	{
	//		typ = val.Type().Elem()
	//		val = reflect.New(typ).Elem()
	//		m.parseStructFields(typ, val)
	//	}
	default:
		log.Panic("object kind [%v] not support yet", typ.Kind())
	}
	return
}

// parse struct fields
func parseStructFields(typ reflect.Type, val reflect.Value) (labels, values []string) {
	kind := typ.Kind()
	if kind == reflect.Struct {
		NumField := val.NumField()
		for i := 0; i < NumField; i++ {
			typField := typ.Field(i)
			valField := val.Field(i)

			if typField.Type.Kind() == reflect.Ptr {
				typField.Type = typField.Type.Elem()
				valField = valField.Elem()
			}
			if !valField.IsValid() || !valField.CanInterface() {
				continue
			}
			strTagVal, ignore := getTagValue(typField)
			if ignore {
				continue
			}
			if strTagVal != "" {
				labels = append(labels, strTagVal)
				values = append(values, fmt.Sprintf("%v", valField.Interface()))
			}
		}
	}
	return
}

// get struct field's tag value
func getTagValue(sf reflect.StructField) (strValue string, ignore bool) {
	strValue = sf.Tag.Get(TagNameLabel)
	if strValue == TagValueIgnore {
		return "", true
	}
	if strValue == "" {
		strValue = sf.Tag.Get(TagNameJson)
		if strValue == TagValueIgnore {
			return "", true
		}
	}
	return
}
